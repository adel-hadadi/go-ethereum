package dto

import (
	"math/big"
	"time"
)

type Transaction struct {
	Amount      *big.Int   `jsonL:"amount"`
	FromAddress string    `json:"from_address"`
	ToAddress   string    `json:"to_address"`
	Time        time.Time `json:"time"`
}
