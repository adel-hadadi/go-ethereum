package service

import (
	"context"
	"os"

	"ethereum/internal/ethereum/adapters"
	"ethereum/internal/ethereum/app"
	"ethereum/internal/ethereum/app/services"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func NewApplication(ctx context.Context) app.Application {
	db, err := sqlx.Connect("sqlite3", "./docker-db/database.db")
	if err != nil {
		panic(err)
	}

	schema := `
CREATE TABLE IF NOT EXISTS transactions (
	id INTEGER primary key autoincrement,
	transaction_id VARCHAR(255) not null unique ,
	from_address VARCHAR(255) not null ,
	to_address VARCHAR(255) not null,
	amount INTEGER UNSIGNED not null
)
`
	db.MustExec(schema)

	ethClient, err := ethclient.Dial(os.Getenv("NETWORK_ADDRESS"))
	if err != nil {
		panic(err)
	}

	ethereumSvc := adapters.NewEthereumService(ethClient)

	transactionRepo := adapters.NewTransactionRepository(db)

	return app.Application{
		Services: app.Services{
			TransactionService: services.NewTransactionService(ethereumSvc, transactionRepo),
			WalletService:      services.NewWalletService(ethereumSvc, transactionRepo),
		},
		Repositories: app.Repositories{
			TransactionRepository: transactionRepo,
		},
	}
}
