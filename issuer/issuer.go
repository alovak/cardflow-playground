package issuer

import (
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
		ID:       uuid.New().String(),
		Balance:  req.Balance,
		Currency: req.Currency,
	}

	err := i.repo.CreateAccount(account)
	if err != nil {
		return nil, fmt.Errorf("creating account: %w", err)
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
		if err == ErrNotFound {
			return AuthorizationResponse{
				ApprovalCode: ApprovalCodeCardInvalid,
			}, nil
		}

		return AuthorizationResponse{}, fmt.Errorf("finding card: %w", err)
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

	// find the account
	// check if the account has enough balance

	// if no, return a insufficient funds error

	// if yes, create transaction for card and return an approval code
	transaction.ApprovalCode = ApprovalCodeApproved
	transaction.AuthorizationCode = "123456"
	transaction.Status = TransactionStatusAuthorized

	return AuthorizationResponse{
		AuthorizationCode: transaction.AuthorizationCode,
		ApprovalCode:      transaction.ApprovalCode,
	}, nil
}

// generateFakeCardNumber generates a fake card number starting with 9
// and a random 15-digit number. This is not a valid card number.
func generateFakeCardNumber() string {
	rand.Seed(time.Now().UnixNano())

	// Generate a 15-digit random number starting with 9
	randomDigits := make([]int, 16)
	randomDigits[0] = 9
	for i := 1; i < len(randomDigits); i++ {
		randomDigits[i] = rand.Intn(10)
	}

	var cardNumber string
	for _, digit := range randomDigits {
		cardNumber += fmt.Sprintf("%d", digit)
	}

	return cardNumber
}
