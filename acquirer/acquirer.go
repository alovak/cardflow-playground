package acquirer

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Acquirer struct {
	repo          Repository
	iso8583Client ISO8583Client
}

type ISO8583Client interface {
	AuthorizePayment(payment *Payment, card Card) (AuthorizationResponse, error)
}

func NewAcquirer(repo Repository, iso8583Client ISO8583Client) *Acquirer {
	return &Acquirer{
		repo:          repo,
		iso8583Client: iso8583Client,
	}
}

func (a *Acquirer) CreateMerchant(create CreateMerchant) (*Merchant, error) {
	merchant := &Merchant{
		ID:         uuid.New().String(),
		Name:       create.Name,
		MCC:        create.MCC,
		PostalCode: create.PostalCode,
		WebSite:    create.WebSite,
	}

	err := a.repo.CreateMerchant(merchant)
	if err != nil {
		return nil, fmt.Errorf("creating merchant: %w", err)
	}

	return merchant, nil
}

func (a *Acquirer) CreatePayment(merchantID string, create CreatePayment) (*Payment, error) {
	payment := &Payment{
		ID:         uuid.New().String(),
		MerchantID: merchantID,
		Amount:     create.Amount,
		Currency:   create.Currency,
		Card: SafeCard{
			First6:         create.Card.Number[:6],
			Last4:          create.Card.Number[len(create.Card.Number)-4:],
			ExpirationDate: create.Card.ExpirationDate,
		},
		Status:    PaymentStatusPending,
		CreatedAt: time.Now(),
	}

	err := a.repo.CreatePayment(payment)
	if err != nil {
		return nil, fmt.Errorf("creating payment: %w", err)
	}

	response, err := a.iso8583Client.AuthorizePayment(payment, create.Card)
	if err != nil {
		payment.Status = PaymentStatusError
		// update payment details
		return nil, fmt.Errorf("authorizing payment: %w", err)
	}

	payment.AuthorizationCode = response.AuthorizationCode

	if response.ApprovalCode == "00" {
		payment.Status = PaymentStatusAuthorized
	} else {
		payment.Status = PaymentStatusDeclined
	}

	return payment, nil
}

func (a *Acquirer) GetPayment(merchantID, paymentID string) (*Payment, error) {
	payment, err := a.repo.GetPayment(merchantID, paymentID)
	if err != nil {
		return nil, fmt.Errorf("getting payment: %w", err)
	}

	return payment, nil
}
