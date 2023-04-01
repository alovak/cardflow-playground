package issuer

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Issuer struct {
	repo Repository
}

func NewIssuer(repo Repository) *Issuer {
	return &Issuer{
		repo: repo,
	}
}

func (i *Issuer) CreateAccount(req CreateAccount) (*Account, error) {
	account := &Account{
		ID:               uuid.New().String(),
		AvailableBalance: req.Balance,
		Currency:         req.Currency,
	}

	err := i.repo.CreateAccount(account)
	if err != nil {
		return nil, fmt.Errorf("creating account: %w", err)
	}

	return account, nil
}

func (i *Issuer) GetAccount(accountID string) (*Account, error) {
	account, err := i.repo.GetAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("finding account: %w", err)
	}

	return account, nil
}

func (i *Issuer) IssueCard(accountID string) (*Card, error) {
	card := &Card{
		ID:                    uuid.New().String(),
		AccountID:             accountID,
		Number:                generateFakeCardNumber(),
		ExpirationDate:        time.Now().AddDate(3, 1, 0).Format("0106"), // 3 years, 1 month from now
		CardVerificationValue: "1234",
	}

	err := i.repo.CreateCard(card)
	if err != nil {
		return nil, fmt.Errorf("creating card: %w", err)
	}

	return card, nil
}

// ListTransactions returns a list of transactions for the given account ID.
func (i *Issuer) ListTransactions(accountID string) ([]*Transaction, error) {
	transactions, err := i.repo.ListTransactions(accountID)
	if err != nil {
		return nil, fmt.Errorf("listing transactions: %w", err)
	}

	return transactions, nil
}

func (i *Issuer) AuthorizeRequest(req AuthorizationRequest) (AuthorizationResponse, error) {
	card, err := i.repo.FindCardForAuthorization(req.Card)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return AuthorizationResponse{
				ApprovalCode: ApprovalCodeInvalidCard,
			}, nil
		}

		return AuthorizationResponse{}, fmt.Errorf("finding card: %w", err)
	}

	account, err := i.repo.GetAccount(card.AccountID)
	if err != nil {
		return AuthorizationResponse{}, fmt.Errorf("finding account: %w", err)
	}

	transaction := &Transaction{
		ID:        uuid.New().String(),
		AccountID: card.AccountID,
		CardID:    card.ID,
		Amount:    req.Amount,
		Currency:  req.Currency,
	}

	err = i.repo.CreateTransaction(transaction)
	if err != nil {
		return AuthorizationResponse{}, fmt.Errorf("creating transaction: %w", err)
	}

	// hold the funds on the account
	err = account.Hold(req.Amount)
	if err != nil {
		// handle insufficient funds
		if !errors.Is(err, ErrInsufficientFunds) {
			return AuthorizationResponse{}, fmt.Errorf("holding funds: %w", err)
		}

		return AuthorizationResponse{
			ApprovalCode: ApprovalCodeInsufficientFunds,
		}, nil
	}

	transaction.ApprovalCode = ApprovalCodeApproved
	transaction.AuthorizationCode = generateAuthorizationCode()
	transaction.Status = TransactionStatusAuthorized

	return AuthorizationResponse{
		AuthorizationCode: transaction.AuthorizationCode,
		ApprovalCode:      transaction.ApprovalCode,
	}, nil
}

// generateFakeCardNumber generates a fake card number starting with 9
// and a random 15-digit number. This is not a valid card number.
func generateFakeCardNumber() string {
	return fmt.Sprintf("9%s", generateRandomNumber(15))
}

func generateAuthorizationCode() string {
	return generateRandomNumber(6)
}

func generateRandomNumber(length int) string {
	rand.Seed(time.Now().UnixNano())

	// Generate a 6-digit random number
	randomDigits := make([]int, length)
	for i := 0; i < len(randomDigits); i++ {
		randomDigits[i] = rand.Intn(10)
	}

	var number string
	for _, digit := range randomDigits {
		number += fmt.Sprintf("%d", digit)
	}

	return number
}
