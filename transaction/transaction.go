package transaction

import (
	"context"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
)

type Transaction struct {
	DB *bun.DB

	tx  *bun.Tx
	ctx context.Context
}

func NewTransaction(db *bun.DB) (*Transaction, error) {
	return &Transaction{DB: db}, nil
}

func (t *Transaction) Init() error {
	ctx := context.Background()
	tx, err := t.DB.BeginTx(ctx, nil)
	t.tx = &tx
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	return nil
}

func (t *Transaction) GetContext() context.Context {
	return t.ctx
}

// Exec Execute a function within the transaction context and if there is no
// errors in it, it will commit the transaction, otherwise it will roll it back
func (t *Transaction) Exec(fn func(context.Context, *bun.Tx) error) error {
	err := fn(t.ctx, t.tx)
	if err != nil {
		t.rollbackTx()
		return errors.New("error executing transaction")
	} else {
		if err := t.commitTx(); err != nil {
			t.rollbackTx()
			return errors.New("error committing transaction")
		}
	}
	return nil
}

func (t *Transaction) commitTx() error {
	if err := t.tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (t *Transaction) rollbackTx() {
	_ = t.tx.Rollback()
}
