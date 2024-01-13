package acquirer

import (
	"fmt"
	"sync"

	"github.com/alovak/cardflow-playground/acquirer/models"
)

var ErrNotFound = fmt.Errorf("not found")

type Repository struct {
	mu sync.RWMutex

	merchants map[string]*models.Merchant
	payments  map[string]*models.Payment
}

func NewRepository() *Repository {
	return &Repository{
		merchants: make(map[string]*models.Merchant),
		payments:  make(map[string]*models.Payment),
	}
}

func (r *Repository) CreateMerchant(merchant *models.Merchant) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.merchants[merchant.ID] = merchant

	return nil
}

func (r *Repository) GetMerchant(merchantID string) (*models.Merchant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	merchant, ok := r.merchants[merchantID]
	if !ok {
		return nil, ErrNotFound
	}

	return merchant, nil
}

func (r *Repository) CreatePayment(payment *models.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.payments[payment.ID] = payment

	return nil
}

func (r *Repository) GetPayment(merchantID, paymentID string) (*models.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	payment, ok := r.payments[paymentID]
	if !ok {
		return nil, ErrNotFound
	}

	if payment.MerchantID != merchantID {
		return nil, ErrNotFound
	}

	return payment, nil
}
