package transaction

import (
	"context"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
)

type Transaction struct {
	ctx context.Context

	DB *bun.DB
}

func NewTransaction(db *bun.DB) (Transaction, error) {
	return Transaction{DB: db}, nil
}

func (t *Transaction) Init() (*bun.Tx, error) {
	ctx := context.Background()
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &tx, nil
}

func (t *Transaction) GetContext() context.Context {
	return t.ctx
}

// Exec Execute a function within the transaction context and if there is no
// errors in it, it will commit the transaction, otherwise it will roll it back
func (t *Transaction) Exec(tx *bun.Tx, fn func(context.Context, *bun.Tx) error) error {
	err := fn(t.ctx, tx)
	if err != nil {
		t.rollbackTx(tx)
		return errors.New("error executing transaction")
	} else {
		if err := t.commitTx(tx); err != nil {
			t.rollbackTx(tx)
			return errors.New("error committing transaction")
		}
	}
	return nil
}

func (t *Transaction) commitTx(tx *bun.Tx) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (t *Transaction) rollbackTx(tx *bun.Tx) {
	_ = tx.Rollback()
}
