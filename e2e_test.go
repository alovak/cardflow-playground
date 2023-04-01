package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	main "github.com/alovak/cardflow-playground"
	"github.com/alovak/cardflow-playground/acquirer"
	"github.com/alovak/cardflow-playground/issuer"
	"github.com/alovak/cardflow-playground/log"
	"github.com/stretchr/testify/require"
)

func TestEndToEndTransaction(t *testing.T) {
	// Initialize the issuer and acquirer components here
	issuerBasePath, iso8583ServerAddr := setupIssuer(t)
	acquirerBasePath := setupAcquirer(t, iso8583ServerAddr)

	// configure the issuer client
	issuerClient := NewIssuerClient(issuerBasePath)
	acquirerClient := NewAcquirerClient(acquirerBasePath)

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
	merchant, err := acquirerClient.CreateMerchant(acquirer.CreateMerchant{
		Name:       "Demo Merchant",
		MCC:        "5411",
		PostalCode: "12345",
		WebSite:    "https://demo.merchant.com",
	})
	require.NoError(t, err)

	// When: Acquirer receives the payment request for the merchant with the issued card
	payment, err := acquirerClient.CreatePayment(merchant.ID, acquirer.CreatePayment{
		Card: acquirer.Card{
			Number:                card.Number,
			CardVerificationValue: card.CardVerificationValue,
			ExpirationDate:        card.ExpirationDate,
		},
		Amount:   10_00, // $10
		Currency: "USD",
	})
	require.NoError(t, err)

	// Then: There should be an authorized transaction in the acquirer
	payment, err = acquirerClient.GetPayment(merchant.ID, payment.ID)
	require.NoError(t, err)
	require.Equal(t, acquirer.PaymentStatusAuthorized, payment.Status)

	// In the issuer, there should be an authorized transaction for the card
	transactions, err := issuerClient.GetTransactions(accountID)
	require.NoError(t, err)

	require.Len(t, transactions, 1)
	require.Equal(t, card.ID, transactions[0].CardID)
	require.Equal(t, 10_00, transactions[0].Amount)
	require.Equal(t, "USD", transactions[0].Currency)
	require.Equal(t, issuer.TransactionStatusAuthorized, transactions[0].Status)
	require.Equal(t, payment.AuthorizationCode, transactions[0].AuthorizationCode)

	// Account's available balance should be less by the transaction amount
	account, err := issuerClient.GetAccount(accountID)
	require.NoError(t, err)

	require.Equal(t, 100_00-10_00, account.AvailableBalance)
	require.Equal(t, 10_00, account.HoldBalance)
}

func setupIssuer(t *testing.T) (string, string) {
	app := issuer.NewApp(log.New())
	err := app.Start()
	require.NoError(t, err)

	// dont' forget to shutdown the issuer app
	t.Cleanup(app.Shutdown)

	return fmt.Sprintf("http://%s", app.Addr), app.ISO8583ServerAddr
}

func setupAcquirer(t *testing.T, iso8583ServerAddr string) string {
	acquirerApp := main.NewAcquirerApp(log.New(), iso8583ServerAddr)
	err := acquirerApp.Start()
	require.NoError(t, err)

	// dont' forget to shutdown the acquirer app
	t.Cleanup(acquirerApp.Shutdown)

	return fmt.Sprintf("http://%s", acquirerApp.Addr)
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

// GetAccount returns the account for the given account ID or an error.
func (i *issuerClient) GetAccount(accountID string) (issuer.Account, error) {
	res, err := i.httpClient.Get(i.baseURL + "/accounts/" + accountID)
	if err != nil {
		return issuer.Account{}, err
	}

	if res.StatusCode != http.StatusOK {
		return issuer.Account{}, fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusOK)
	}

	var account issuer.Account
	err = json.NewDecoder(res.Body).Decode(&account)
	if err != nil {
		return issuer.Account{}, err
	}

	return account, nil
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

// GetTransactions returns the list of transactions for the given card ID
// and account ID or an error.
func (i *issuerClient) GetTransactions(accountID string) ([]issuer.Transaction, error) {
	res, err := i.httpClient.Get(i.baseURL + "/accounts/" + accountID + "/transactions")
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusOK)
	}

	var transactions []issuer.Transaction
	err = json.NewDecoder(res.Body).Decode(&transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

type acquirerClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewAcquirerClient(baseURL string) *acquirerClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: 5 * time.Second,
		},
	}

	return &acquirerClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *acquirerClient) CreateMerchant(req acquirer.CreateMerchant) (acquirer.Merchant, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return acquirer.Merchant{}, err
	}

	res, err := c.httpClient.Post(c.baseURL+"/merchants", "application/json", bytes.NewReader(reqJSON))
	if err != nil {
		return acquirer.Merchant{}, err
	}

	if res.StatusCode != http.StatusCreated {
		return acquirer.Merchant{}, fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusCreated)
	}

	var merchant acquirer.Merchant
	err = json.NewDecoder(res.Body).Decode(&merchant)
	if err != nil {
		return acquirer.Merchant{}, err
	}

	return merchant, nil
}

func (c *acquirerClient) CreatePayment(merchantID string, req acquirer.CreatePayment) (acquirer.Payment, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return acquirer.Payment{}, err
	}

	res, err := c.httpClient.Post(c.baseURL+"/merchants/"+merchantID+"/payments", "application/json", bytes.NewReader(reqJSON))
	if err != nil {
		return acquirer.Payment{}, err
	}

	if res.StatusCode != http.StatusCreated {
		return acquirer.Payment{}, fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusCreated)
	}

	var payment acquirer.Payment
	err = json.NewDecoder(res.Body).Decode(&payment)
	if err != nil {
		return acquirer.Payment{}, err
	}

	return payment, nil
}

func (c *acquirerClient) GetPayment(merchantID, paymentID string) (acquirer.Payment, error) {
	res, err := c.httpClient.Get(c.baseURL + "/merchants/" + merchantID + "/payments/" + paymentID)
	if err != nil {
		return acquirer.Payment{}, err
	}

	if res.StatusCode != http.StatusOK {
		return acquirer.Payment{}, fmt.Errorf("unexpected status code: %d; expected: %d", res.StatusCode, http.StatusOK)
	}

	var payment acquirer.Payment
	err = json.NewDecoder(res.Body).Decode(&payment)
	if err != nil {
		return acquirer.Payment{}, err
	}

	return payment, nil
}
