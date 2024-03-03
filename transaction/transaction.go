package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
)

type Transaction struct {
	ctx context.Context
	tx  *bun.Tx
}

func NewTransaction(db *bun.DB) (*Transaction, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &Transaction{tx: &tx, ctx: ctx}, nil
}

func (t *Transaction) GetContext() context.Context {
	return t.ctx
}

// Exec Execute a function within the transaction context
func (t *Transaction) Exec(fn func(context.Context, *bun.Tx) error) error {
	err := fn(t.ctx, t.tx)
	if err != nil {
		t.RollbackTx()
		return errors.New("failed to execute function within transaction: " + err.Error())
	}
	return nil
}

func (t *Transaction) CommitTx() error {
	if err := t.tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (t *Transaction) RollbackTx() {
	_ = t.tx.Rollback()
}
