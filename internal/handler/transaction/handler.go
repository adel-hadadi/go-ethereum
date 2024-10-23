package transactionhandler

import (
	"context"
	"errors"
	"ethereum/util/respond"
	"math/big"
	"net/http"
	"regexp"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type (
	Handler struct {
		ethereumService ethereumService
	}

	ethereumService interface {
		CreateTransaction(ctx context.Context, toAddress string, amount int64) (*types.Transaction, error)
		CheckTransaction(ctx context.Context, txID string, amount *big.Int) error
	}

	CreateTransactionReq struct {
		Value     int64  `json:"value"`
		ToAddress string `json:"to_address"`
	}

	CreateTransactionRes struct {
		TransactionID string `json:"transaction_id"`
	}
)

func (t *CreateTransactionReq) Bind(r *http.Request) error {
	if t.ToAddress == "" {
		return errors.New("to_address is required")
	}

	if t.Value == 0 {
		return errors.New("value is required")
	}

	if ok, _ := regexp.Match("^(0x)?[0-9a-fA-F]{40}$", []byte(t.ToAddress)); !ok {
		return errors.New("to_address is invalid!")
	}

	return nil
}

func New(service ethereumService) Handler {
	return Handler{
		ethereumService: service,
	}
}

func (h Handler) Routes(r chi.Router) {
	r.Post("/transactions", h.Create)
	r.Post("/transactions/check", h.CheckStatus)
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionReq
	if err := render.Bind(r, &req); err != nil {
		respond.Faile(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	tx, err := h.ethereumService.CreateTransaction(r.Context(), req.ToAddress, req.Value)
	if err != nil {
		respond.WithErr(w, err)
		return
	}

	respond.Success(w, r, CreateTransactionRes{
		TransactionID: tx.Hash().String(),
	}, http.StatusCreated)
}
