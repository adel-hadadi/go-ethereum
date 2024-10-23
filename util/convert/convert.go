package convert

import (
	"math"
	"math/big"
)

func ToETH(balance *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())

	return new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
}
