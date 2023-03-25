package issuer

import (
	"fmt"

	"github.com/google/uuid"
)

type Issuer struct {
	repo Repository
}

func New(repo Repository) *Issuer {
	return &Issuer{
		repo: repo,
	}
}

func (i *Issuer) CreateAccount(req CreateAccountRequest) (*Account, error) {
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
