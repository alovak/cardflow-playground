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

func (a *API) AppendRoutes(router chi.Router) {
	router.Post("/accounts", a.createAccount)
}

func (a *API) createAccount(w http.ResponseWriter, r *http.Request) {
	create := CreateAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(&create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, err := a.issuer.CreateAccount(create)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}
