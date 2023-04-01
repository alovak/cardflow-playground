package issuer

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type API struct {
	issuer *Service
}

func NewAPI(issuer *Service) *API {
	return &API{
		issuer: issuer,
	}
}

func (a *API) AppendRoutes(r chi.Router) {
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", a.createAccount)
		r.Route("/{accountID}", func(r chi.Router) {
			r.Get("/", a.getAccount)
			r.Post("/cards", a.issueCard)
			r.Get("/transactions", a.getTransactions)
		})
	})
}

func (a *API) createAccount(w http.ResponseWriter, r *http.Request) {
	create := CreateAccount{}
	err := json.NewDecoder(r.Body).Decode(&create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, err := a.issuer.CreateAccount(create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (a *API) getAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")

	account, err := a.issuer.GetAccount(accountID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func (a *API) issueCard(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")

	card, err := a.issuer.IssueCard(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

func (a *API) getTransactions(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")

	transactions, err := a.issuer.ListTransactions(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}
