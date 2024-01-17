package models

type Transaction struct {
	ID                string
	AccountID         string
	CardID            string
	Amount            int64
	Currency          string
	AuthorizationCode string
	ApprovalCode      string
	Status            TransactionStatus
	Merchant          Merchant
}

type TransactionStatus string

const (
	TransactionStatusAuthorized TransactionStatus = "authorized"
	TransactionStatusDeclined   TransactionStatus = "declined"
)
