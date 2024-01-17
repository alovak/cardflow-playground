package iso8583

import (
	"fmt"

	"github.com/alovak/cardflow-playground/issuer/models"
	"github.com/moov-io/iso8583"
	iso8583Connection "github.com/moov-io/iso8583-connection"
	iso8583Server "github.com/moov-io/iso8583-connection/server"
	"golang.org/x/exp/slog"
)

// Server is a wrapper around the moov-io/iso8583-connection server.
type Server struct {
	Addr string

	server     *iso8583Server.Server
	logger     *slog.Logger
	authorizer Authorizer
}

// Authorizer is an interface that defines the authorization logic.
type Authorizer interface {
	AuthorizeRequest(req models.AuthorizationRequest) (models.AuthorizationResponse, error)
}

// NewServer creates a new Server instance with the given logger, address and authorizer.
func NewServer(logger *slog.Logger, addr string, authorizer Authorizer) *Server {
	logger = logger.With(slog.String("type", "iso8583-server"), slog.String("addr", addr))

	s := &Server{
		logger:     logger,
		Addr:       addr,
		authorizer: authorizer,
	}

	// here we create an instance of the ISO 8583 server
	iso8583Server := iso8583Server.New(
		// this is the ISO 8583 spec we defined in spec.go
		spec,

		// part of binary framing, it reads the message length from the connection
		readMessageLength,

		// part of binary framing, it writes the message length to the connection
		writeMessageLength,

		// here we define a function that will be called when a new message is received`
		iso8583Connection.InboundMessageHandler(s.handleRequest),
	)

	s.server = iso8583Server

	return s
}

// Start starts the server.
func (s *Server) Start() error {
	s.logger.Info("starting ISO 8583 server...")

	if err := s.server.Start(s.Addr); err != nil {
		return fmt.Errorf("starting ISO 8583 server: %w", err)
	}

	// if the server is started successfully, update the address as it might be
	// different from the one we passed to the Start() method (e.g. if we passed
	// ":0" to let the OS choose a free port)
	s.Addr = s.server.Addr

	s.logger.Info("ISO 8583 server started", slog.String("addr", s.Addr))

	return nil
}

func (s *Server) Close() error {
	s.logger.Info("shutting down ISO 8583 server...")

	s.server.Close()

	s.logger.Info("ISO 8583 server shut down")

	return nil
}

// handleRequest is called when a new message is received.
func (s *Server) handleRequest(c *iso8583Connection.Connection, message *iso8583.Message) {
	mti, err := message.GetMTI()
	if err != nil {
		s.logger.Error("failed to get MTI from message", "err", err)
	}

	logger := s.logger.With(slog.String("mti", mti))

	logger.Info("handling request")

	// here we handle different MTIs
	switch mti {
	case "0100":
		err = s.handleAuthorizationRequest(c, message)
	default:
		err = fmt.Errorf("unknown MTI: %s", mti)
	}

	if err != nil {
		logger.Error("failed to handle request", "err", err)
	}
}

// handleAuthorizationRequest handles authorization requests.
func (s *Server) handleAuthorizationRequest(c *iso8583Connection.Connection, message *iso8583.Message) error {
	// here we unmarshal the message into our AuthorizationRequest struct
	requestData := &AuthorizationRequest{}
	if err := message.Unmarshal(requestData); err != nil {
		return fmt.Errorf("unmarshaling message: %w", err)
	}

	s.logger.With(
		slog.String("mti", requestData.MTI),
		slog.String("stan", requestData.STAN),
		slog.Int64("amount", requestData.Amount),
		slog.String("currency", requestData.Currency),
	).Info("handling authorization request")

	// here we create an instance of our authorization request
	// and pass it to the authorizer
	authRequest := models.AuthorizationRequest{
		Amount:   requestData.Amount,
		Currency: requestData.Currency,
		Card: models.Card{
			Number:                requestData.PrimaryAccountNumber,
			ExpirationDate:        requestData.ExpirationDate,
			CardVerificationValue: requestData.CardVerificationValue,
		},
		Merchant: models.Merchant{
			Name:       requestData.AcceptorInformation.Name,
			MCC:        requestData.AcceptorInformation.MCC,
			PostalCode: requestData.AcceptorInformation.PostalCode,
			WebSite:    requestData.AcceptorInformation.WebSite,
		},
	}

	// we define a variable that will hold the response data
	// we need to define it here so we can set its value in the if/else block
	var responseData *AuthorizationResponse

	// pass the request to the authorizer and get the response with the
	// approval code and authorization code
	authResponse, err := s.authorizer.AuthorizeRequest(authRequest)
	if err != nil {
		responseData = &AuthorizationResponse{
			MTI:          "0110",
			STAN:         requestData.STAN,
			ApprovalCode: models.ApprovalCodeSystemError,
		}
	} else {
		responseData = &AuthorizationResponse{
			MTI:               "0110",
			STAN:              requestData.STAN,
			ApprovalCode:      authResponse.ApprovalCode,
			AuthorizationCode: authResponse.AuthorizationCode,
		}
	}

	// create response message and marshal the response data into it
	responseMessage := iso8583.NewMessage(spec)
	if err := responseMessage.Marshal(responseData); err != nil {
		return fmt.Errorf("marshaling response: %w", err)
	}

	// send the response message back to the client
	if err := c.Reply(responseMessage); err != nil {
		return fmt.Errorf("sending response: %w", err)
	}

	s.logger.With(
		slog.String("mti", responseData.MTI),
		slog.String("stan", responseData.STAN),
		slog.String("approval_code", responseData.ApprovalCode),
		slog.String("authorization_code", responseData.AuthorizationCode),
	).Info("authorization response sent")

	return nil
}
