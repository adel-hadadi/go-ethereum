package ports

import (
	"errors"
	"math/big"
	"net/http"
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

type TransactionRes struct {
	ID          int        `json:"id"`
	TxID        string     `json:"txid"`
	FromAddress string     `json:"from_address"`
	ToAddress   string     `json:"to_address"`
	Amount      *big.Float `json:"amount"`
}

type CreateTransactionReq struct {
	ToAddress string   `json:"to_address"`
	Amount    *big.Int `json:"amount"`
}

func (t *CreateTransactionReq) Bind(r *http.Request) error {
	if t.ToAddress == "" {
		return errors.New("to_address is required")
	}

	if t.Amount == nil {
		return errors.New("amount is required")
	}

	return nil
}
