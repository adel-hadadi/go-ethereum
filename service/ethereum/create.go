package ethereumservice

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"ethereum/util/apperr"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func (s EthereumService) CreateTransaction(ctx context.Context, toAddress string, amount int64) (*types.Transaction, error) {
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

	nonce, err := s.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get pending nonce: %w", err))
	}

	value := big.NewInt(amount)
	gasLimit := uint64(21000)

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on get suggested gas price: %w", err))
	}

	toAddressHex := common.HexToAddress(toAddress)

	tx := types.NewTransaction(nonce, toAddressHex, value, gasLimit, gasPrice, nil)

	chainID, err := s.client.NetworkID(ctx)
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
		return nil, apperr.New(apperr.Unexpected).WithErr(fmt.Errorf("Error on decode raw transaction: %w", err))
	}

	tx = new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)
	s.client.SendTransaction(ctx, tx)

	return tx, nil
}
