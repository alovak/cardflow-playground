package issuer

import "sync"

type Repository interface {
	CreateAccount(account *Account) error
}

type repository struct {
	Cards    []*Card
	Accounts []*Account

	mu sync.RWMutex
}

func NewRepository() *repository {
	return &repository{
		Cards:    []*Card{},
		Accounts: []*Account{},
	}
}

func (r *repository) CreateAccount(account *Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Accounts = append(r.Accounts, account)

	return nil
}
