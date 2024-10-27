package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BenjaminArmijo3/bank/internal/db/sqlc"
)

type Store struct {
	*sqlc.Queries
	db *sql.DB
}

func NewStore(database *sql.DB) *Store {
	return &Store{
		db:      database,
		Queries: sqlc.New(database),
	}
}

func (s *Store) ExecTx(ctx context.Context, fq func(q *sqlc.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
	err = fq(q)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("encountered rollback error: %v", txErr)
		}
		return err
	}

	return tx.Commit()
}
