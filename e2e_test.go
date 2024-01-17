package main_test

import (
	"fmt"
	"testing"

	"github.com/alovak/cardflow-playground/acquirer"
	acquirerClient "github.com/alovak/cardflow-playground/acquirer/client"
	"github.com/alovak/cardflow-playground/acquirer/models"
	"github.com/alovak/cardflow-playground/issuer"
	issuerClient "github.com/alovak/cardflow-playground/issuer/client"
	issuerModels "github.com/alovak/cardflow-playground/issuer/models"
	"github.com/alovak/cardflow-playground/log"
	"github.com/stretchr/testify/require"
)

func TestEndToEndTransaction(t *testing.T) {
	// Initialize the issuer and acquirer components here
	issuerBasePath, iso8583ServerAddr := setupIssuer(t)
	acquirerBasePath := setupAcquirer(t, iso8583ServerAddr)

	// configure the issuer client
	issuerClient := issuerClient.New(issuerBasePath)
	acquirerClient := acquirerClient.New(acquirerBasePath)

	// Given: Create an account with $100 balance
	accountID, err := issuerClient.CreateAccount(issuerModels.CreateAccount{
		Balance:  100_00, // $100
		Currency: "USD",
	})
	require.NoError(t, err)

	// Issue a card for the account
	card, err := issuerClient.IssueCard(accountID)
	require.NoError(t, err)

	// Given: Create a new merchant for the acquirer
	merchant, err := acquirerClient.CreateMerchant(models.CreateMerchant{
		Name:       "Demo Merchant",
		MCC:        "5411",
		PostalCode: "12345",
		WebSite:    "https://demo.merchant.com",
	})
	require.NoError(t, err)

	// When: Acquirer receives the payment request for the merchant with the issued card
	payment, err := acquirerClient.CreatePayment(merchant.ID, models.CreatePayment{
		Card: models.Card{
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
	require.Equal(t, models.PaymentStatusAuthorized, payment.Status)

	// In the issuer, there should be an authorized transaction for the card
	transactions, err := issuerClient.GetTransactions(accountID)
	require.NoError(t, err)

	require.Len(t, transactions, 1)
	require.Equal(t, card.ID, transactions[0].CardID)

	// check the transaction details
	require.Equal(t, int64(10_00), transactions[0].Amount)
	require.Equal(t, "USD", transactions[0].Currency)
	require.Equal(t, issuerModels.TransactionStatusAuthorized, transactions[0].Status)
	require.Equal(t, payment.AuthorizationCode, transactions[0].AuthorizationCode)

	// check the merchant details
	require.Equal(t, merchant.Name, transactions[0].Merchant.Name)
	require.Equal(t, merchant.MCC, transactions[0].Merchant.MCC)
	require.Equal(t, merchant.PostalCode, transactions[0].Merchant.PostalCode)
	require.Equal(t, merchant.WebSite, transactions[0].Merchant.WebSite)

	// Account's available balance should be less by the transaction amount
	account, err := issuerClient.GetAccount(accountID)
	require.NoError(t, err)

	require.Equal(t, int64(100_00-10_00), account.AvailableBalance)
	require.Equal(t, int64(10_00), account.HoldBalance)
}

func setupIssuer(t *testing.T) (string, string) {
	app := issuer.NewApp(log.New(), &issuer.Config{
		HTTPAddr:    "127.0.0.1:0", // use random port
		ISO8583Addr: "127.0.0.1:0", // use random port
	})
	err := app.Start()
	require.NoError(t, err)

	// dont' forget to shutdown the issuer app
	t.Cleanup(app.Shutdown)

	return fmt.Sprintf("http://%s", app.Addr), app.ISO8583ServerAddr
}

func setupAcquirer(t *testing.T, iso8583ServerAddr string) string {
	app := acquirer.NewApp(log.New(), &acquirer.Config{
		HTTPAddr:    "127.0.0.1:0", // use random port
		ISO8583Addr: iso8583ServerAddr,
	})
	err := app.Start()
	require.NoError(t, err)

	// dont' forget to shutdown the acquirer app
	t.Cleanup(app.Shutdown)

	return fmt.Sprintf("http://%s", app.Addr)
}
