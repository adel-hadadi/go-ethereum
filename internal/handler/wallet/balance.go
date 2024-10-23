package wallethandler

import (
	"ethereum/util/convert"
	"ethereum/util/respond"
	"math/big"
	"net/http"
)

type BalanceRes struct {
	Balance *big.Float `json:"balance"`
}

func (h Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	balance, err := h.ethereumSvc.GetBalance(r.Context())
	if err != nil {
		respond.WithErr(w, err)
		return
	}

	res := BalanceRes{
		Balance: convert.ToETH(balance),
	}

	respond.Success(w, r, res)
}
