package ports

import (
	"math/big"
	"net/http"

	"ethereum/internal/ethereum/app"
	"ethereum/internal/ethereum/app/services"
	"ethereum/util/convert"
	"ethereum/util/respond"

	"github.com/go-chi/render"
)

type HttpServer struct {
	app app.Application
}

func NewHTTPServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}

// WalletBalance TODO: move to proper location
type WalletBalance struct {
	Balance *big.Float `json:"balance"`
}

func (h HttpServer) Balance(w http.ResponseWriter, r *http.Request) {
	balance, err := h.app.Services.WalletService.GetBalance(r.Context())
	if err != nil {
		respond.WithErr(w, err)
		return
	}

	respond.Success(w, r, WalletBalance{convert.ToETH(balance)})
}

func (h HttpServer) Deposit(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.app.Services.WalletService.Deposit(r.Context())
	if err != nil {
		respond.WithErr(w, err)
		return
	}

	res := transactionToResponse(transactions)

	respond.Success(w, r, res)
}

func (h HttpServer) Withdraw(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.app.Services.WalletService.Withdraw(r.Context())
	if err != nil {
		respond.WithErr(w, err)
		return
	}

	res := transactionToResponse(transactions)

	respond.Success(w, r, res)
}

func transactionToResponse(transactions []services.Transaction) []TransactionRes {
	res := make([]TransactionRes, 0, len(transactions))

	for _, t := range transactions {
		res = append(res, TransactionRes{
			ID:          t.ID,
			TxID:        t.TransactionID,
			FromAddress: t.FromAddress,
			ToAddress:   t.ToAddress,
			Amount:      convert.ToETH(t.Amount),
		})
	}

	return res
}

func (h HttpServer) TransactionCheck(w http.ResponseWriter, r *http.Request) {
	var req TransactionCheckReq
	if err := render.Bind(r, &req); err != nil {
		respond.Faile(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	err := h.app.Services.TransactionService.Check(r.Context(), req.TxID, req.Amount)
	if err != nil {
		respond.WithErr(w, err)
		return
	}

	respond.Success(w, r, nil)
}

func (h HttpServer) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionReq
	if err := render.Bind(r, &req); err != nil {
		respond.Faile(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	transaction, err := h.app.Services.TransactionService.Create(r.Context(), req.Amount, req.ToAddress)
	if err != nil {
		respond.WithErr(w, err)
		return
	}

	respond.Success(w, r, transaction, http.StatusCreated)
}
