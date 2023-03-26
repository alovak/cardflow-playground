package acquirer

import (
	"fmt"
	"sync"
)

var ErrNotFound = fmt.Errorf("not found")

type Repository interface {
	CreateMerchant(merchant *Merchant) error
	CreatePayment(payment *Payment) error
	GetPayment(merchantID, paymentID string) (*Payment, error)
}

type repository struct {
	mu sync.RWMutex

	merchants map[string]*Merchant
	payments  map[string]*Payment
}

func NewRepository() *repository {
	return &repository{
		merchants: make(map[string]*Merchant),
		payments:  make(map[string]*Payment),
	}
}

func (r *repository) CreateMerchant(merchant *Merchant) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.merchants[merchant.ID] = merchant

	return nil
}

func (r *repository) CreatePayment(payment *Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.payments[payment.ID] = payment

	return nil
}

func (r *repository) GetPayment(merchantID, paymentID string) (*Payment, error) {
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
