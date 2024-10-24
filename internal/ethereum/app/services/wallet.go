package services

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"ethereum/internal/ethereum/adapters"
	"ethereum/util/apperr"
)

type WalletService struct {
	ethereumSvc     ethereumService
	transactionRepo transactionRepository
}

type Transaction struct {
	ID            int
	TransactionID string
	FromAddress   string
	ToAddress     string
	Amount        *big.Int
}

func NewWalletService(
	ethereumSvc ethereumService,
	transactionRepo transactionRepository,
) WalletService {
	return WalletService{
		ethereumSvc:     ethereumSvc,
		transactionRepo: transactionRepo,
	}
}

func (s WalletService) GetBalance(ctx context.Context) (*big.Int, error) {
	balance, err := s.ethereumSvc.GetWalletBalance(ctx, os.Getenv("ACCOUNT_ADDRESS"))
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).
			WithErr(fmt.Errorf("Error on get account balance: %w", err)).
			WithMessage("cant get account balance. try later")
	}

	return balance, nil
}

func (s WalletService) Deposit(ctx context.Context) ([]Transaction, error) {
	transactions, err := s.transactionRepo.GetWalletTransactions(
		ctx,
		adapters.GetWalletTransactionsFilter{To: os.Getenv("ACCOUNT_ADDRESS")},
	)
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get wallet deposit: %w", err))
	}

	transactionRes := make([]Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		transactionRes = append(transactionRes, Transaction{
			ID:            transaction.ID,
			TransactionID: transaction.TransactionID,
			FromAddress:   transaction.FromAddress,
			ToAddress:     transaction.ToAddress,
			Amount:        big.NewInt(transaction.Amount),
		})
	}

	return transactionRes, nil
}

func (s WalletService) Withdraw(ctx context.Context) ([]Transaction, error) {
	transactions, err := s.transactionRepo.GetWalletTransactions(
		ctx,
		adapters.GetWalletTransactionsFilter{
			From: os.Getenv("ACCOUNT_ADDRESS"),
		},
	)
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get wallet deposit: %w", err))
	}

	transactionRes := make([]Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		transactionRes = append(transactionRes, Transaction{
			ID:            transaction.ID,
			TransactionID: transaction.TransactionID,
			FromAddress:   transaction.FromAddress,
			ToAddress:     transaction.ToAddress,
			Amount:        big.NewInt(transaction.Amount),
		})
	}

	return transactionRes, nil
}
