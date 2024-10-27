package store

import (
	"context"

	"github.com/BenjaminArmijo3/bank/internal/db/sqlc"
)

type TransferTxResponse struct {
	FromAccount sqlc.Account  `json:"from_account"`
	ToAccount   sqlc.Account  `json:"to_account"`
	Transfer    sqlc.Transfer `json:"transfer"`
}

func (s *Store) TransferTx(ctx context.Context, tr sqlc.CreateTransferParams) (TransferTxResponse, error) {
	var tx TransferTxResponse
	var errT error

	err := s.ExecTx(ctx, func(q *sqlc.Queries) error {
		// transfer money
		tArg := sqlc.CreateTransferParams{
			FromAccountID: tr.FromAccountID, ToAccountID: tr.ToAccountID, Amount: tr.Amount,
		}
		tx.Transfer, errT = q.CreateTransfer(context.Background(), tArg)
		if errT != nil {
			return errT
		}

		toArg := sqlc.UpdateAccountBalanceNewParams{
			Amount: tr.Amount,
			ID:     int64(tr.ToAccountID),
		}
		tx.ToAccount, errT = s.Queries.UpdateAccountBalanceNew(context.Background(), toArg)
		if errT != nil {
			return errT
		}
		fromArg := sqlc.UpdateAccountBalanceNewParams{
			Amount: tr.Amount * -1,
			ID:     int64(tr.FromAccountID),
		}
		tx.FromAccount, errT = q.UpdateAccountBalanceNew(context.Background(), fromArg)
		if errT != nil {
			return errT
		}
		return nil
	})
	return tx, err
}
