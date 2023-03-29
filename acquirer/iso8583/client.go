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
	return c.iso8583Connection.Connect()
}

func (c *Client) AuthorizePayment(payment *acquirer.Payment) (acquirer.AuthorizationResponse, error) {
	return acquirer.AuthorizationResponse{}, nil
}
