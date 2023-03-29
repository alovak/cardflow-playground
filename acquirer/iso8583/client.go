package iso8583

import (
	"fmt"

	"github.com/alovak/cardflow-playground/acquirer"
	iso8583Connection "github.com/moov-io/iso8583-connection"
)

type Client struct {
	iso8583Connection *iso8583Connection.Connection
}

func NewClient(iso8583ServerAddr string) (*Client, error) {
	conn, err := iso8583Connection.New(iso8583ServerAddr, spec, readMessageLength, writeMessageLength)
	if err != nil {
		return nil, fmt.Errorf("creating iso8583 connection: %w", err)
	}

	return &Client{
		iso8583Connection: conn,
	}, nil
}

func (c *Client) Connect() error {
	fmt.Printf("Connecting to ISO 8583 Server %s...\n", c.iso8583Connection.Addr())

	if err := c.iso8583Connection.Connect(); err != nil {
		return fmt.Errorf("connecting to ISO 8583 Server: %w", err)
	}

	fmt.Printf("Connected to ISO 8583 Server %s\n", c.iso8583Connection.Addr())
	return nil
}

func (c *Client) AuthorizePayment(payment *acquirer.Payment) (acquirer.AuthorizationResponse, error) {
	return acquirer.AuthorizationResponse{}, nil
}
