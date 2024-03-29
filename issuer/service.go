package issuer

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/alovak/cardflow-playground/issuer/models"
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (i *Service) CreateAccount(req models.CreateAccount) (*models.Account, error) {
	account := &models.Account{
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

func (i *Service) GetAccount(accountID string) (*models.Account, error) {
	account, err := i.repo.GetAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("finding account: %w", err)
	}

	return account, nil
}

func (i *Service) IssueCard(accountID string) (*models.Card, error) {
	card := &models.Card{
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
func (i *Service) ListTransactions(accountID string) ([]*models.Transaction, error) {
	transactions, err := i.repo.ListTransactions(accountID)
	if err != nil {
		return nil, fmt.Errorf("listing transactions: %w", err)
	}

	return transactions, nil
}

func (i *Service) AuthorizeRequest(req models.AuthorizationRequest) (models.AuthorizationResponse, error) {
	card, err := i.repo.FindCardForAuthorization(req.Card)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return models.AuthorizationResponse{
				ApprovalCode: models.ApprovalCodeInvalidCard,
			}, nil
		}

		return models.AuthorizationResponse{}, fmt.Errorf("finding card: %w", err)
	}

	account, err := i.repo.GetAccount(card.AccountID)
	if err != nil {
		return models.AuthorizationResponse{}, fmt.Errorf("finding account: %w", err)
	}

	transaction := &models.Transaction{
		ID:        uuid.New().String(),
		AccountID: card.AccountID,
		CardID:    card.ID,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Merchant:  req.Merchant,
	}

	err = i.repo.CreateTransaction(transaction)
	if err != nil {
		return models.AuthorizationResponse{}, fmt.Errorf("creating transaction: %w", err)
	}

	// hold the funds on the account
	err = account.Hold(req.Amount)
	if err != nil {
		// handle insufficient funds
		if !errors.Is(err, models.ErrInsufficientFunds) {
			return models.AuthorizationResponse{}, fmt.Errorf("holding funds: %w", err)
		}

		return models.AuthorizationResponse{
			ApprovalCode: models.ApprovalCodeInsufficientFunds,
		}, nil
	}

	transaction.ApprovalCode = models.ApprovalCodeApproved
	transaction.AuthorizationCode = generateAuthorizationCode()
	transaction.Status = models.TransactionStatusAuthorized

	return models.AuthorizationResponse{
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
