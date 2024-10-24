package services

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"ethereum/util/apperr"
	"github.com/ethereum/go-ethereum"
)

type TransactionService struct {
	ethereumSvc     ethereumService
	transactionRepo transactionRepository
}

func NewTransactionService(
	ethereumSvc ethereumService,
	transactionRepo transactionRepository,
) TransactionService {
	return TransactionService{
		ethereumSvc:     ethereumSvc,
		transactionRepo: transactionRepo,
	}
}

func (s TransactionService) Check(ctx context.Context, txID string, amount *big.Int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.ethereumSvc.GetTransaction(ctx, txID)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return apperr.New(apperr.NotFound).WithMessage("Transaction with given id not found")
		}

		return apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get transcation by hash: %w", err))
	}

	if tx.Value().String() != amount.String() {
		return apperr.New(apperr.BadRequest).WithMessage("Transaction is invalid")
	}

	if tx.To().Hex() != os.Getenv("ACCOUNT_ADDRESS") {
		return apperr.New(apperr.BadRequest).WithMessage("Transaction is not for us")
	}

	sender, err := s.ethereumSvc.GetSender(ctx, tx)
	if err != nil {
		return apperr.New(apperr.Unexpected).
			WithErr(fmt.Errorf("Error on find sender: %w", err))
	}

	walletAddr := os.Getenv("ACCOUNT_ADDRESS")

	err = s.transactionRepo.Create(ctx, tx.Hash().Hex(), sender, walletAddr, tx.Value(), tx.Time())
	if err != nil && !apperr.IsSQLDuplicateEntry(err) {
		return apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on inserting transactions: %w", err))
	}

	return nil
}

func (s TransactionService) Create(ctx context.Context, amount *big.Int, to string) (Transaction, error) {
	from := os.Getenv("ACCOUNT_ADDRESS")

	transaction, err := s.ethereumSvc.CreateTransaction(ctx, amount, from, to)
	if err != nil {
		return Transaction{}, err
	}

	err = s.transactionRepo.Create(ctx, transaction.Hash().Hex(), from, to, amount, transaction.Time())
	if err != nil {
		return Transaction{}, apperr.New(apperr.Unexpected).
			WithErr(fmt.Errorf("error on inserting transaction to database: %w", err))
	}

	return Transaction{
		ID:            0,
		TransactionID: transaction.Hash().Hex(),
		FromAddress:   from,
		ToAddress:     to,
		Amount:        amount,
	}, nil
}
