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
		ID:             uuid.New().String(),
		AccountID:      accountID,
		Number:         generateFakeCardNumber(),
		ExpirationDate: time.Now().AddDate(3, 1, 0).Format("01/06"), // 3 years, 1 month from now
		CVV:            "123",
	}

	err := i.repo.CreateCard(card)
	if err != nil {
		return nil, fmt.Errorf("creating card: %w", err)
	}

	return card, nil
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
