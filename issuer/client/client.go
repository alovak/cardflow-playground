package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alovak/cardflow-playground/issuer/models"
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

// CreateAccount creates a new account with the given balance and currency and
// returns the account ID or an error.
func (i *client) CreateAccount(req models.CreateAccount) (string, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	res, err := i.httpClient.Post(i.baseURL+"/accounts", "application/json", bytes.NewReader(reqJSON))
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusCreated)
	}

	var account models.Account
	err = json.NewDecoder(res.Body).Decode(&account)
	if err != nil {
		return "", err
	}

	return account.ID, nil
}

// GetAccount returns the account for the given account ID or an error.
func (i *client) GetAccount(accountID string) (models.Account, error) {
	res, err := i.httpClient.Get(i.baseURL + "/accounts/" + accountID)
	if err != nil {
		return models.Account{}, err
	}

	if res.StatusCode != http.StatusOK {
		return models.Account{}, fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusOK)
	}

	var account models.Account
	err = json.NewDecoder(res.Body).Decode(&account)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

// IssueCard issues a new card for the given account ID and returns the card or
// an error.
func (i *client) IssueCard(accountID string) (models.Card, error) {
	res, err := i.httpClient.Post(i.baseURL+"/accounts/"+accountID+"/cards", "application/json", nil)
	if err != nil {
		return models.Card{}, err
	}

	if res.StatusCode != http.StatusCreated {
		return models.Card{}, fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusCreated)
	}

	var card models.Card
	err = json.NewDecoder(res.Body).Decode(&card)
	if err != nil {
		return models.Card{}, err
	}

	return card, nil
}

// GetTransactions returns the list of transactions for the given card ID
// and account ID or an error.
func (i *client) GetTransactions(accountID string) ([]models.Transaction, error) {
	res, err := i.httpClient.Get(i.baseURL + "/accounts/" + accountID + "/transactions")
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusOK)
	}

	var transactions []models.Transaction
	err = json.NewDecoder(res.Body).Decode(&transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
