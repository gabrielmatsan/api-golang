package helpers

import (
	"context"

	"github.com/gabrielmatsan/teste-api/cmd/db"
	"github.com/jmoiron/sqlx"
)

func BeginTransaction(ctx context.Context, db *sqlx.DB) (*sqlx.Tx, error) {
	return db.BeginTxx(ctx, nil)
}

func WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	db := db.GetDB()
	tx, err := BeginTransaction(ctx, db)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return fn(tx)
}
