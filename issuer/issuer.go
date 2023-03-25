package issuer

import "github.com/google/uuid"

type Issuer struct {
	Accounts []*Account
	Cards    []Card
}

func New() *Issuer {
	return &Issuer{
		Accounts: make([]*Account, 0),
		Cards:    make([]Card, 0),
	}
}

func (i *Issuer) CreateAccount(req CreateAccountRequest) (*Account, error) {
	account := &Account{
		ID:       uuid.New().String(),
		Balance:  req.Balance,
		Currency: req.Currency,
	}

	i.Accounts = append(i.Accounts, account)

	return account, nil
}
