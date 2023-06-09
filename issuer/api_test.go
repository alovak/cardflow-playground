package issuer_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alovak/cardflow-playground/issuer"
	"github.com/alovak/cardflow-playground/issuer/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

// Just a simple test
func TestAPI(t *testing.T) {
	router := chi.NewRouter()

	api := issuer.NewAPI(issuer.NewService(issuer.NewRepository()))
	api.AppendRoutes(router)

	t.Run("create account", func(t *testing.T) {
		create := models.CreateAccount{
			Balance:  10_00,
			Currency: "USD",
		}

		jsonReq, _ := json.Marshal(create)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonReq))
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		account := models.Account{}
		err := json.Unmarshal(w.Body.Bytes(), &account)
		require.NoError(t, err)

		require.Equal(t, create.Balance, account.AvailableBalance)
		require.Equal(t, create.Currency, account.Currency)
		require.NotEmpty(t, account.ID)
	})
}
