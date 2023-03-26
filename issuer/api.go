package issuer

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type API struct {
	issuer *Issuer
}

func NewAPI(issuer *Issuer) *API {
	return &API{
		issuer: issuer,
	}
}

func (a *API) AppendRoutes(r chi.Router) {
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", a.createAccount)
		r.Route("/{accountID}", func(r chi.Router) {
			r.Post("/cards", a.issueCard)
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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (a *API) issueCard(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")

	card, err := a.issuer.IssueCard(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}
