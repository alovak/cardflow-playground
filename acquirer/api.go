package acquirer

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type API struct {
	acquirer *Acquirer
}

func NewAPI(acquirer *Acquirer) *API {
	return &API{
		acquirer: acquirer,
	}
}

func (a *API) AppendRoutes(r chi.Router) {
	r.Route("/merchants", func(r chi.Router) {
		r.Post("/", a.createMerchant)
		r.Route("/{merchantID}", func(r chi.Router) {
			r.Post("/payments", a.createPayment)
			r.Get("/payments/{paymentID}", a.getPayment)
		})
	})
}

func (a *API) createMerchant(w http.ResponseWriter, r *http.Request) {
	create := CreateMerchant{}
	err := json.NewDecoder(r.Body).Decode(&create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, err := a.acquirer.CreateMerchant(create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (a *API) createPayment(w http.ResponseWriter, r *http.Request) {
	merchantID := chi.URLParam(r, "merchantID")

	create := CreatePayment{}
	err := json.NewDecoder(r.Body).Decode(&create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment, err := a.acquirer.CreatePayment(merchantID, create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

func (a *API) getPayment(w http.ResponseWriter, r *http.Request) {
	merchantID := chi.URLParam(r, "merchantID")
	paymentID := chi.URLParam(r, "paymentID")

	payment, err := a.acquirer.GetPayment(merchantID, paymentID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}