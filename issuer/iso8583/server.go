package iso8583

import (
	"fmt"
	"strconv"

	"github.com/alovak/cardflow-playground/issuer/models"
	"github.com/moov-io/iso8583"
	iso8583Connection "github.com/moov-io/iso8583-connection"
	iso8583Server "github.com/moov-io/iso8583-connection/server"
	"github.com/moov-io/iso8583/field"
	"golang.org/x/exp/slog"
)

type server struct {
	Addr string

	server     *iso8583Server.Server
	logger     *slog.Logger
	authorizer Authorizer
}

type Authorizer interface {
	AuthorizeRequest(req models.AuthorizationRequest) (models.AuthorizationResponse, error)
}

func NewServer(logger *slog.Logger, addr string, authorizer Authorizer) *server {
	logger = logger.With(slog.String("type", "iso8583-server"), slog.String("addr", addr))

	s := &server{
		logger:     logger,
		Addr:       addr,
		authorizer: authorizer,
	}

	iso8583Server := iso8583Server.New(
		spec,
		readMessageLength,
		writeMessageLength,
		iso8583Connection.InboundMessageHandler(s.handleRequest),
	)

	s.server = iso8583Server

	return s
}

func (s *server) Start() error {
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

func (s *server) Close() error {
	s.logger.Info("shutting down ISO 8583 server...")

	s.server.Close()

	s.logger.Info("ISO 8583 server shut down")

	return nil
}

func (s *server) handleRequest(c *iso8583Connection.Connection, message *iso8583.Message) {
	mti, err := message.GetMTI()
	if err != nil {
		s.logger.Error("failed to get MTI from message", "err", err)
	}

	logger := s.logger.With(slog.String("mti", mti))

	logger.Info("handling request")

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

func (s *server) handleAuthorizationRequest(c *iso8583Connection.Connection, message *iso8583.Message) error {
	requestData := &AuthorizationRequest{}
	if err := message.Unmarshal(requestData); err != nil {
		return fmt.Errorf("unmarshaling message: %w", err)
	}

	s.logger.With(
		slog.String("mti", requestData.MTI.Value()),
		slog.String("stan", requestData.STAN.Value()),
		slog.String("amount", requestData.Amount.Value()),
		slog.String("currency", requestData.Currency.Value()),
	).Info("handling authorization request")

	amount, err := strconv.Atoi(requestData.Amount.Value())
	if err != nil {
		return fmt.Errorf("parsing amount: %w", err)
	}

	authRequest := models.AuthorizationRequest{
		Amount:   amount,
		Currency: requestData.Currency.Value(),
		Card: models.Card{
			Number:                requestData.PrimaryAccountNumber.Value(),
			ExpirationDate:        requestData.ExpirationDate.Value(),
			CardVerificationValue: requestData.CardVerificationValue.Value(),
		},
		Merchant: models.Merchant{
			Name:       requestData.AcceptorInformation.Name.Value(),
			MCC:        requestData.AcceptorInformation.MCC.Value(),
			PostalCode: requestData.AcceptorInformation.PostalCode.Value(),
			WebSite:    requestData.AcceptorInformation.WebSite.Value(),
		},
	}

	var responseData *AuthorizationResponse

	authResponse, err := s.authorizer.AuthorizeRequest(authRequest)
	if err != nil {
		return fmt.Errorf("authorizing request: %w", err)

		responseData = &AuthorizationResponse{
			MTI:          field.NewStringValue("0110"),
			STAN:         field.NewStringValue(requestData.STAN.Value()),
			ApprovalCode: field.NewStringValue(models.ApprovalCodeSystemError),
		}
	} else {
		responseData = &AuthorizationResponse{
			MTI:               field.NewStringValue("0110"),
			STAN:              field.NewStringValue(requestData.STAN.Value()),
			ApprovalCode:      field.NewStringValue(authResponse.ApprovalCode),
			AuthorizationCode: field.NewStringValue(authResponse.AuthorizationCode),
		}
	}

	responseMessage := iso8583.NewMessage(spec)
	if err := responseMessage.Marshal(responseData); err != nil {
		return fmt.Errorf("marshaling response: %w", err)
	}

	if err := c.Reply(responseMessage); err != nil {
		return fmt.Errorf("sending response: %w", err)
	}

	s.logger.With(
		slog.String("mti", responseData.MTI.Value()),
		slog.String("stan", responseData.STAN.Value()),
		slog.String("approval_code", responseData.ApprovalCode.Value()),
		slog.String("authorization_code", responseData.AuthorizationCode.Value()),
	).Info("authorization response sent")

	return nil
}
