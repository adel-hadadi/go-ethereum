package ethereumservice

import (
	"context"
	"errors"
	"ethereum/data/dto"
	"ethereum/util/apperr"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func (s EthereumService) GetBalance(ctx context.Context) (*big.Int, error) {
	account := common.HexToAddress(os.Getenv("ACCOUNT_ADDRESS"))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	balance, err := s.client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).
			WithErr(fmt.Errorf("Error on get account balance: %w", err)).
			WithMessage("cant get account balance. try later")
	}

	return balance, nil
}

func (s EthereumService) CheckTransaction(ctx context.Context, txID string, amount *big.Int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	txHash := common.HexToHash(txID)

	tx, _, err := s.client.TransactionByHash(ctx, txHash)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return apperr.New(apperr.NotFound).WithMessage("Transaction with given id not found")
		}

		return apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get transcation by hash: %w", err))
	}

	if tx.Value().String() != amount.String() {
		return apperr.New(apperr.BadRequest).WithMessage("Transaction is invalid")
	}

	return nil
}

func (s EthereumService) GetTransactionsByWallet(ctx context.Context) ([]dto.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	walletAddr := common.HexToAddress(os.Getenv("ACCOUNT_ADDRESS"))

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   nil,
		Addresses: []common.Address{walletAddr},
	}

	logs, err := s.client.FilterLogs(ctx, query)
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get logs: %w", err))
	}

	transactions := make([]dto.Transaction, 0)

	for _, v := range logs {
		_, err := s.client.BlockByHash(ctx, v.BlockHash)
		if err != nil {
			return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get block by block hash: %w", err))
		}

		log.Println("transaction id => ", v.TxHash)
	}

	// startBlock := big.NewInt(0)
	// endBlock := "latest"

	// latestBlock, err := s.client.BlockNumber(ctx)
	// if err != nil {
	// 	return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get latest block: %w", err))
	// }

	// if endBlock == "latest" {
	// 	endBlockInt := new(big.Int).SetUint64(latestBlock)
	// 	endBlock = endBlockInt.String()
	// }

	// for i := startBlock.Int64(); i < int64(latestBlock); i++ {
	// 	block, err := s.client.BlockByNumber(ctx, big.NewInt(i))
	// 	if err != nil {
	// 		return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get block by number: %w", err))
	// 	}

	// 	// TODO: check block time
	// 	for _, tx := range block.Transactions() {
	// 		if tx.To() != nil && tx.To().Hex() == walletAddr.Hex() {
	// 			chainID, err := s.client.NetworkID(ctx)
	// 			if err != nil {
	// 				return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get networkd id: %w", err))
	// 			}

	// 			from, err := types.Sender(types.NewLondonSigner(chainID), tx)
	// 			if err != nil {
	// 				return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on find from address in transaction: %w", err))
	// 			}

	// 			transactions = append(transactions, dto.Transaction{
	// 				Amount:      tx.Value(),
	// 				FromAddress: from.Hex(),
	// 			})
	// 		}
	// 	}
	// }

	return transactions, nil
}
