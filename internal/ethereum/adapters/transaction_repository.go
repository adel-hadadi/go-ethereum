package adapters

import (
	"context"
	"math/big"
	"time"

	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	db *sqlx.DB
}

type TransactionModel struct {
	ID            int    `json:"id" db:"id"`
	FromAddress   string `json:"from_address" db:"from_address"`
	ToAddress     string `json:"to_address" db:"to_address"`
	Amount        int64  `json:"amount" db:"amount"`
	TransactionID string `json:"transaction_id" db:"transaction_id"`
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return TransactionRepository{db: db}
}

func (r TransactionRepository) Create(ctx context.Context, txID, from, to string, amount *big.Int, _ time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO transactions (transaction_id, from_address, to_address, amount) VALUES($1, $2, $3, $4)
	`

	floatAmount, _ := amount.Float64()
	_, err := r.db.ExecContext(ctx, query, txID, from, to, floatAmount)
	if err != nil {
		return err
	}

	return nil
}

type GetWalletTransactionsFilter struct {
	From string
	To   string
}

func (r TransactionRepository) GetWalletTransactions(ctx context.Context, filters GetWalletTransactionsFilter) ([]TransactionModel, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, transaction_id, from_address, to_address, amount FROM transactions 
		WHERE (
			to_address=$1 OR $1 = ''
		) AND (
		    from_address=$2 OR $2 = ''
		)
`

	rows, err := r.db.QueryxContext(ctx, query, filters.To, filters.From)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []TransactionModel
	for rows.Next() {
		t := TransactionModel{}

		if err := rows.Scan(
			&t.ID,
			&t.TransactionID,
			&t.FromAddress,
			&t.ToAddress,
			&t.Amount,
		); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}
