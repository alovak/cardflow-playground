package models

type Transaction struct {
	ID                string
	AccountID         string
	CardID            string
	Amount            int
	Currency          string
	AuthorizationCode string
	ApprovalCode      string
	Status            TransactionStatus
}

type TransactionStatus string

const (
	TransactionStatusAuthorized TransactionStatus = "authorized"
	TransactionStatusDeclined   TransactionStatus = "declined"
)
