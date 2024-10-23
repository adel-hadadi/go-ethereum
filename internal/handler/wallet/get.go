package wallethandler

import (
	"ethereum/util/respond"
	"net/http"
)

func (h Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.ethereumSvc.GetTransactionsByWallet(r.Context())
	if err != nil {
		respond.WithErr(w, err)
		return
	}

	respond.Success(w, r, transactions)
}
