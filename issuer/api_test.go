package issuer_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alovak/cardflow-playground/issuer"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	router := chi.NewRouter()

	api := issuer.NewAPI(issuer.New())
	api.AppendRoutes(router)

	t.Run("create account", func(t *testing.T) {
		create := issuer.CreateAccountRequest{
			Balance:  1000,
			Currency: "USD",
		}

		jsonReq, _ := json.Marshal(create)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonReq))
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		account := issuer.Account{}
		err := json.Unmarshal(w.Body.Bytes(), &account)
		require.NoError(t, err)

		require.Equal(t, create.Balance, account.Balance)
		require.Equal(t, create.Currency, account.Currency)
		require.NotEmpty(t, account.ID)
	})
}
