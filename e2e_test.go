package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alovak/cardflow-playground/issuer"
	"github.com/stretchr/testify/require"
)

func TestEndToEndTransaction(t *testing.T) {
	// Initialize the issuer and acquirer components here
	issuerBasePath := setupIssuer(t)

	// configure the issuer client
	issuerClient := NewIssuerClient(issuerBasePath)

	// Given: Create an account with $100 balance
	accountID, err := issuerClient.CreateAccount(issuer.CreateAccount{
		Balance:  100_00, // $100
		Currency: "USD",
	})
	require.NoError(t, err)

	// Issue a card for the account
	card, err := issuerClient.IssueCard(accountID)
	require.NoError(t, err)

	// Given: Create a new merchant for the acquirer

	// When: Acquirer receives the payment request for the merchant with the issued card

	// Then: There should be an authorized transaction in the acquirer

	// In the issuer, there should be an authorized transaction for the card

	// Account's available balance should be less by the transaction amount

	// Account's hold balance should be equal to the transaction amount
}

func setupIssuer(t *testing.T) string {
	issuerApp := issuer.NewApp()
	err := issuerApp.Start()
	require.NoError(t, err)

	// dont' forget to shutdown the issuer app
	t.Cleanup(issuerApp.Shutdown)

	return fmt.Sprintf("http://%s", issuerApp.Addr)
}

type issuerClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewIssuerClient(baseURL string) *issuerClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: 5 * time.Second,
		},
	}

	return &issuerClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// CreateAccount creates a new account with the given balance and currency and
// returns the account ID or an error.
func (i *issuerClient) CreateAccount(req issuer.CreateAccount) (string, error) {
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

	var account issuer.Account
	err = json.NewDecoder(res.Body).Decode(&account)
	if err != nil {
		return "", err
	}

	return account.ID, nil
}

// IssueCard issues a new card for the given account ID and returns the card or
// an error.
func (i *issuerClient) IssueCard(accountID string) (issuer.Card, error) {
	res, err := i.httpClient.Post(i.baseURL+"/accounts/"+accountID+"/cards", "application/json", nil)
	if err != nil {
		return issuer.Card{}, err
	}

	if res.StatusCode != http.StatusCreated {
		return issuer.Card{}, fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusCreated)
	}

	var card issuer.Card
	err = json.NewDecoder(res.Body).Decode(&card)
	if err != nil {
		return issuer.Card{}, err
	}

	return card, nil
}
