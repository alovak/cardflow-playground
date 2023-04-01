package issuer

import (
	"fmt"
	"sync"
)

var ErrNotFound = fmt.Errorf("not found")

type Repository interface {
	CreateAccount(account *Account) error
	GetAccount(accountID string) (*Account, error)
	CreateCard(card *Card) error
	ListTransactions(accountID string) ([]*Transaction, error)
	FindCardForAuthorization(card Card) (*Card, error)
	CreateTransaction(transaction *Transaction) error
}

type repository struct {
	Cards        []*Card
	Accounts     []*Account
	Transactions []*Transaction

	mu sync.RWMutex
}

func NewRepository() *repository {
	return &repository{
		Cards:        []*Card{},
		Accounts:     []*Account{},
		Transactions: []*Transaction{},
	}
}

func (r *repository) CreateAccount(account *Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Accounts = append(r.Accounts, account)

	return nil
}

func (r *repository) GetAccount(accountID string) (*Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, account := range r.Accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return nil, ErrNotFound
}

func (r *repository) CreateCard(card *Card) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Cards = append(r.Cards, card)

	return nil
}

func (r *repository) FindCardForAuthorization(card Card) (*Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, c := range r.Cards {
		match := c.Number == card.Number &&
			c.ExpirationDate == card.ExpirationDate &&
			c.CardVerificationValue == card.CardVerificationValue

		if match {
			return c, nil
		}
	}

	return nil, ErrNotFound
}

func (r *repository) CreateTransaction(transaction *Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Transactions = append(r.Transactions, transaction)

	return nil
}

// ListTransactions returns all transactions for a given account ID.
func (r *repository) ListTransactions(accountID string) ([]*Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*Transaction

	for _, transaction := range r.Transactions {
		if transaction.AccountID == accountID {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}
