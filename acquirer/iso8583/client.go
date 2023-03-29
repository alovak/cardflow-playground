package iso8583

import (
	"fmt"

	"github.com/alovak/cardflow-playground/acquirer"
	iso8583Connection "github.com/moov-io/iso8583-connection"
	"golang.org/x/exp/slog"
)

type Client struct {
	iso8583Connection *iso8583Connection.Connection
	logger            *slog.Logger
}

func NewClient(logger *slog.Logger, iso8583ServerAddr string) (*Client, error) {
	logger = logger.With(slog.String("type", "iso8583-client"), slog.String("addr", iso8583ServerAddr))

	conn, err := iso8583Connection.New(iso8583ServerAddr, spec, readMessageLength, writeMessageLength)
	if err != nil {
		return nil, fmt.Errorf("creating iso8583 connection: %w", err)
	}

	return &Client{
		iso8583Connection: conn,
		logger:            logger,
	}, nil
}

func (c *Client) Connect() error {
	c.logger.Info("connecting to ISO 8583 server...")

	if err := c.iso8583Connection.Connect(); err != nil {
		return fmt.Errorf("connecting to ISO 8583 server: %w", err)
	}

	c.logger.Info("connected to ISO 8583 server")
	return nil
}

func (c *Client) AuthorizePayment(payment *acquirer.Payment) (acquirer.AuthorizationResponse, error) {
	return acquirer.AuthorizationResponse{}, nil
}
