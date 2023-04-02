package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alovak/cardflow-playground/acquirer/models"
)

type client struct {
	httpClient *http.Client
	baseURL    string
}

func New(baseURL string) *client {
	httpClient := &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: 5 * time.Second,
		},
	}

	return &client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *client) CreateMerchant(req models.CreateMerchant) (models.Merchant, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return models.Merchant{}, err
	}

	res, err := c.httpClient.Post(c.baseURL+"/merchants", "application/json", bytes.NewReader(reqJSON))
	if err != nil {
		return models.Merchant{}, err
	}

	if res.StatusCode != http.StatusCreated {
		return models.Merchant{}, fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusCreated)
	}

	var merchant models.Merchant
	err = json.NewDecoder(res.Body).Decode(&merchant)
	if err != nil {
		return models.Merchant{}, err
	}

	return merchant, nil
}

func (c *client) CreatePayment(merchantID string, req models.CreatePayment) (models.Payment, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return models.Payment{}, err
	}

	res, err := c.httpClient.Post(c.baseURL+"/merchants/"+merchantID+"/payments", "application/json", bytes.NewReader(reqJSON))
	if err != nil {
		return models.Payment{}, err
	}

	if res.StatusCode != http.StatusCreated {
		return models.Payment{}, fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusCreated)
	}

	var payment models.Payment
	err = json.NewDecoder(res.Body).Decode(&payment)
	if err != nil {
		return models.Payment{}, err
	}

	return payment, nil
}

func (c *client) GetPayment(merchantID, paymentID string) (models.Payment, error) {
	res, err := c.httpClient.Get(c.baseURL + "/merchants/" + merchantID + "/payments/" + paymentID)
	if err != nil {
		return models.Payment{}, err
	}

	if res.StatusCode != http.StatusOK {
		return models.Payment{}, fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusOK)
	}

	var payment models.Payment
	err = json.NewDecoder(res.Body).Decode(&payment)
	if err != nil {
		return models.Payment{}, err
	}

	return payment, nil
}
