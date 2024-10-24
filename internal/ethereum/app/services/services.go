package services

import (
	"context"
	"math/big"
	"time"

	"ethereum/internal/ethereum/adapters"

	"github.com/ethereum/go-ethereum/core/types"
)

type ethereumService interface {
	GetWalletBalance(ctx context.Context, walletAddress string) (*big.Int, error)
	GetTransaction(ctx context.Context, txID string) (*types.Transaction, error)
	GetSender(ctx context.Context, tx *types.Transaction) (string, error)
	CreateTransaction(ctx context.Context, amount *big.Int, from string, to string) (*types.Transaction, error)
}

type transactionRepository interface {
	Create(ctx context.Context, txID, from, to string, amount *big.Int, createdAt time.Time) error
	GetWalletTransactions(ctx context.Context, filters adapters.GetWalletTransactionsFilter) ([]adapters.TransactionModel, error)
}
