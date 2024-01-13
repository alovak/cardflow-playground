package acquirer

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alovak/cardflow-playground/acquirer/models"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

type API struct {
	acquirer *Service
	logger   *slog.Logger
}

func NewAPI(logger *slog.Logger, acquirer *Service) *API {
	return &API{
		logger:   logger,
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
	create := models.CreateMerchant{}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (a *API) createPayment(w http.ResponseWriter, r *http.Request) {
	merchantID := chi.URLParam(r, "merchantID")

	create := models.CreatePayment{}
	err := json.NewDecoder(r.Body).Decode(&create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment, err := a.acquirer.CreatePayment(merchantID, create)
	if err != nil {
		a.logger.Error("failed to create payment", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}
