package issuer

import (
	"fmt"
	"sync"

	"github.com/alovak/cardflow-playground/issuer/models"
)

var ErrNotFound = fmt.Errorf("not found")

type Repository struct {
	Cards        []*models.Card
	Accounts     []*models.Account
	Transactions []*models.Transaction

	mu sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		Cards:        make([]*models.Card, 0),
		Accounts:     make([]*models.Account, 0),
		Transactions: make([]*models.Transaction, 0),
	}
}

func (r *Repository) CreateAccount(account *models.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Accounts = append(r.Accounts, account)

	return nil
}

func (r *Repository) GetAccount(accountID string) (*models.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, account := range r.Accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return nil, ErrNotFound
}

func (r *Repository) CreateCard(card *models.Card) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Cards = append(r.Cards, card)

	return nil
}

func (r *Repository) FindCardForAuthorization(card models.Card) (*models.Card, error) {
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

func (r *Repository) CreateTransaction(transaction *models.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Transactions = append(r.Transactions, transaction)

	return nil
}

// ListTransactions returns all transactions for a given account ID.
func (r *Repository) ListTransactions(accountID string) ([]*models.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*models.Transaction

	for _, transaction := range r.Transactions {
		if transaction.AccountID == accountID {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}
