package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionManager provides a safe abstraction for running PostgreSQL transactions.
type TransactionManager struct {
	pool *pgxpool.Pool
}

// NewTransactionManager creates a new TransactionManager from a pgx connection pool.
func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{pool: pool}
}

// WithTransaction runs the given function inside a database transaction.
// It automatically commits if fn returns nil, or rolls back if an error occurs.
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Ensure rollback if something goes wrong
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p) // rethrow panic
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
