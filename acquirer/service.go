package acquirer

import (
	"fmt"
	"time"

	"github.com/alovak/cardflow-playground/acquirer/models"
	"github.com/google/uuid"
)

type Service struct {
	repo          *Repository
	iso8583Client ISO8583Client
}

type ISO8583Client interface {
	AuthorizePayment(payment *models.Payment, card models.Card, merchant models.Merchant) (models.AuthorizationResponse, error)
}

func NewService(repo *Repository, iso8583Client ISO8583Client) *Service {
	return &Service{
		repo:          repo,
		iso8583Client: iso8583Client,
	}
}

func (a *Service) CreateMerchant(create models.CreateMerchant) (*models.Merchant, error) {
	merchant := &models.Merchant{
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

func (a *Service) CreatePayment(merchantID string, create models.CreatePayment) (*models.Payment, error) {
	payment := &models.Payment{
		ID:         uuid.New().String(),
		MerchantID: merchantID,
		Amount:     create.Amount,
		Currency:   create.Currency,
		Card: models.SafeCard{
			First6:         create.Card.Number[:6],
			Last4:          create.Card.Number[len(create.Card.Number)-4:],
			ExpirationDate: create.Card.ExpirationDate,
		},
		Status:    models.PaymentStatusPending,
		CreatedAt: time.Now(),
	}

	err := a.repo.CreatePayment(payment)
	if err != nil {
		return nil, fmt.Errorf("creating payment: %w", err)
	}

	merchant, err := a.repo.GetMerchant(merchantID)
	if err != nil {
		return nil, fmt.Errorf("getting merchant: %w", err)
	}

	response, err := a.iso8583Client.AuthorizePayment(payment, create.Card, *merchant)
	if err != nil {
		payment.Status = models.PaymentStatusError
		// update payment details
		return nil, fmt.Errorf("authorizing payment: %w", err)
	}

	payment.AuthorizationCode = response.AuthorizationCode

	if response.ApprovalCode == "00" {
		payment.Status = models.PaymentStatusAuthorized
	} else {
		payment.Status = models.PaymentStatusDeclined
	}

	return payment, nil
}

func (a *Service) GetPayment(merchantID, paymentID string) (*models.Payment, error) {
	payment, err := a.repo.GetPayment(merchantID, paymentID)
	if err != nil {
		return nil, fmt.Errorf("getting payment: %w", err)
	}

	return payment, nil
}
