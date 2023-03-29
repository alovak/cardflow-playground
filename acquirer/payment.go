package acquirer

import "time"

type CreatePayment struct {
	Amount   int
	Currency string
	Card     Card
}

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusError      PaymentStatus = "error"
	PaymentStatusAuthorized PaymentStatus = "authorized"
	PaymentStatusDeclined   PaymentStatus = "declined"
)

type Payment struct {
	ID         string
	MerchantID string
	Amount     int
	Currency   string
	Card       SafeCard
	Status     PaymentStatus
	CreatedAt  time.Time
}
