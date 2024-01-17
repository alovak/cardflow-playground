package models

import (
	"errors"
	"sync"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type CreateAccount struct {
	Balance  int64
	Currency string
}

type Account struct {
	ID               string
	AvailableBalance int64
	HoldBalance      int64
	Currency         string

	mu sync.Mutex
}

func (a *Account) Hold(amount int64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.AvailableBalance < amount {
		return ErrInsufficientFunds
	}

	a.AvailableBalance -= amount
	a.HoldBalance += amount

	return nil
}
