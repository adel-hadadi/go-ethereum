package app

import (
	"ethereum/internal/ethereum/adapters"
	"ethereum/internal/ethereum/app/services"
)

type Application struct {
	Services Services

	Repositories Repositories
}

type Services struct {
	TransactionService services.TransactionService
	WalletService      services.WalletService
}

type Repositories struct {
	TransactionRepository adapters.TransactionRepository
}
