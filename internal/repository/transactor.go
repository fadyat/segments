package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type ContextKey string

const (
	ContextTransaction ContextKey = "transaction"
)

// Transactor is an interface for sharing transactions between repository functions.
// It is implemented by the repository package and can be used in the underlying repository functions.
type Transactor interface {

	// RunTx runs the given function in a transaction, rolling back if an error is returned.
	//
	// The function is passed a context with a transaction value, which should be used for all database
	// operations.
	//
	// Usage example:
	//
	// 	err := r.RunTx(context.Background(), nil, func(ctx context.Context) error {
	//		// do repository operations, using ctx to get the transaction
	// 		return nil
	// 	})
	//
	RunTx(context.Context, *sql.TxOptions, func(context.Context) error) error

	// UseTx runs the given function with the existing transaction or using sqlx.DB if transaction is nil.
	//
	// Usage example:
	//
	// err := r.RunTx(context.Background(), nil, func(ctx context.Context) error {
	// 		r.UseTx(ctx, func(executor repository.Executor) {
	// 			// do repository operations, using executor to execute queries
	// 		})
	// 		return nil
	// 	})
	//
	// or can be used in the underlying repository functions:
	//
	// func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Segment, error) {
	// 		r.UseTx(ctx, func(executor repository.Executor) {
	// 			// do repository operations, using executor to execute queries
	// 		})
	// }
	//
	UseTx(context.Context, func(Executor))
}

// Executor is an interface for executing queries, used
// as an interface for the sqlx.DB type and the sqlx.Tx type.
type Executor interface {
	sqlx.ExtContext

	// SelectContext within a transaction/database and context.
	// Any placeholder parameters are replaced with supplied args.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// GetContext within a transaction/database and context.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type transactor struct {
	db *sqlx.DB
}

func NewTransactor(db *sqlx.DB) Transactor {
	return &transactor{db: db}
}

func (t *transactor) RunTx(
	ctx context.Context, level *sql.TxOptions, f func(context.Context) error,
) error {
	tx, err := t.db.BeginTxx(ctx, level)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if e := tx.Rollback(); e != nil {
			err = fmt.Errorf("failed to rollback transaction: %w", e)
		}
	}()

	ctx = context.WithValue(ctx, ContextTransaction, tx)
	if e := f(ctx); e != nil {
		return e
	}

	if e := tx.Commit(); e != nil {
		zap.L().Info("failed to commit transaction", zap.Error(e))
		return e
	}

	return nil
}

func (t *transactor) UseTx(ctx context.Context, f func(Executor)) {
	tx, ok := ctx.Value(ContextTransaction).(*sqlx.Tx)
	if !ok {
		f(t.db)
		return
	}

	f(tx)
}
