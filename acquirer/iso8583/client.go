package iso8583

import (
	"github.com/alovak/cardflow-playground/acquirer"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) AuthorizePayment(payment *acquirer.Payment) (acquirer.AuthorizationResponse, error) {
	return acquirer.AuthorizationResponse{}, nil
}
