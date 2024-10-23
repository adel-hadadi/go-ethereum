package wallethandler

import (
	"context"
	"ethereum/data/dto"
	"math/big"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	ethereumSvc ethereumService
}

type ethereumService interface {
	GetBalance(ctx context.Context) (*big.Int, error)
	GetTransactionsByWallet(ctx context.Context) ([]dto.Transaction, error)
}

func New(service ethereumService) Handler {
	return Handler{
		ethereumSvc: service,
	}
}

func (h Handler) Routes(r chi.Router) {
	r.Get("/wallet/balance", h.GetBalance)
	r.Get("/wallet/transactions", h.GetTransactions)
}
