package transactionhandler

import (
	"errors"
	"ethereum/util/respond"
	"math/big"
	"net/http"

	"github.com/go-chi/render"
)

type TransactionCheckReq struct {
	TxID   string   `json:"txid"`
	Amount *big.Int `json:"amount"`
}

func (t TransactionCheckReq) Bind(r *http.Request) error {
	if t.TxID == "" {
		return errors.New("txid is required")
	}

	if t.Amount == nil {
		return errors.New("amount is required")
	}

	return nil
}

func (h Handler) CheckStatus(w http.ResponseWriter, r *http.Request) {
	var req TransactionCheckReq
	if err := render.Bind(r, &req); err != nil {
		respond.Faile(w, err.Error(), nil)
		return
	}

	err := h.ethereumService.CheckTransaction(r.Context(), req.TxID, req.Amount)
	if err != nil {
		respond.WithErr(w, err)
		return
	}

	respond.Success(w, r, nil, http.StatusOK)
}
