package adapters

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"time"

	"ethereum/util/apperr"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

type EthClient struct {
	client *ethclient.Client
}

func NewEthereumService(client *ethclient.Client) EthClient {
	return EthClient{
		client: client,
	}
}

func (c EthClient) GetWalletBalance(ctx context.Context, walletAddress string) (*big.Int, error) {
	account := common.HexToAddress(walletAddress)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	balance, err := c.client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (c EthClient) GetTransaction(ctx context.Context, txID string) (*types.Transaction, error) {
	txHash := common.HexToHash(txID)

	tx, _, err := c.client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (c EthClient) GetSender(ctx context.Context, tx *types.Transaction) (string, error) {
	chainID, err := c.client.NetworkID(ctx)
	if err != nil {
		return "", err
	}

	msg, err := types.Sender(types.NewLondonSigner(chainID), tx)
	if err != nil {
		return "", err
	}

	return msg.Hex(), nil
}

func (c EthClient) CreateTransaction(ctx context.Context, amount *big.Int, from string, to string) (*types.Transaction, error) {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).
			WithMessage("Cant find private key").
			WithErr(err)
	}

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, apperr.New(apperr.Unexpected).
			WithErr(fmt.Errorf("Error on type assertion for public key: %w", err))
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := c.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get pending nonce: %w", err))
	}

	gasLimit := uint64(21000)

	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get suggested gas price: %w", err))
	}

	toAddressHex := common.HexToAddress(from)

	tx := types.NewTransaction(nonce, toAddressHex, amount, gasLimit, gasPrice, nil)

	chainID, err := c.client.NetworkID(ctx)
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get network id: %w", err))
	}

	signedTX, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on sign transaction: %w", err))
	}

	ts := types.Transactions{signedTX}

	b := new(bytes.Buffer)
	ts.EncodeIndex(0, b)

	rawTxBytes := b.Bytes()

	rawTxHex := hex.EncodeToString(rawTxBytes)

	rawTxBytes, err = hex.DecodeString(rawTxHex)
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).
			WithErr(fmt.Errorf("Error on decode raw transaction: %w", err))
	}

	tx = new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)

	if err := c.client.SendTransaction(ctx, tx); err != nil {
		return nil, apperr.New(apperr.Unexpected).
			WithErr(fmt.Errorf("Error on sending transaction: %w", err))
	}

	return tx, nil
}
