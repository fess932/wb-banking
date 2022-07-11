package account

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
	"net/http"
)

type IUsecase interface {
	Register(ctx context.Context, account *Account) error
	ShowAmount(ctx context.Context, accountID string) (decimal.Decimal, error)
	Transfer(ctx context.Context, fromAccountID, toAccountID string, amount decimal.Decimal) error
}

type HTTPDelivery struct {
	uc IUsecase
}

func NewHTTPDelivery(uc IUsecase) *HTTPDelivery {
	return &HTTPDelivery{uc}
}

type RegisterRequest struct {
	Email  string `json:"email"`
	Amount string `json:"amount"`
}

func (h *HTTPDelivery) Register(w http.ResponseWriter, r *http.Request) {
	rr := &RegisterRequest{}

	if err := json.NewDecoder(r.Body).Decode(rr); err != nil {
		jsonError(w, http.StatusBadRequest, err)
		return
	}

	if rr.Email == "" {
		jsonError(w, http.StatusBadRequest, errors.New("пустой email"))
		return
	}

	amount, err := decimal.NewFromString(rr.Amount)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err)
		return
	}

	acc := &Account{
		Email:  rr.Email,
		Amount: amount,
	}
	if err := h.uc.Register(r.Context(), acc); err != nil {
		jsonError(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, acc)
}

func (h *HTTPDelivery) Amount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		jsonError(w, http.StatusBadRequest, errors.New("emtpy id"))

		return
	}

	amount, err := h.uc.ShowAmount(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err)

		return
	}

	jsonResponse(w, amount)
}

type TransferRequest struct {
	FromID, ToID string
	Amount       string
}

func (h *HTTPDelivery) Transfer(w http.ResponseWriter, r *http.Request) {
	tr := &TransferRequest{}

	if err := json.NewDecoder(r.Body).Decode(tr); err != nil {
		jsonError(w, http.StatusBadRequest, err)

		return
	}

	amount, err := decimal.NewFromString(tr.Amount)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err)

		return
	}

	if err = h.uc.Transfer(r.Context(), tr.FromID, tr.ToID, amount); err != nil {
		jsonError(w, http.StatusInternalServerError, err)

		return
	}

	jsonResponse(w, "ok")
}

func jsonError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func jsonResponse(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"body": body,
	})
}
