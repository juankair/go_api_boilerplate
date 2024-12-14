package dbcontext

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/uptrace/bunrouter"
	"net/http"
)

type DB struct {
	db *dbx.DB
}

type TransactionFunc func(ctx context.Context, f func(ctx context.Context) error) error

type contextKey int

const (
	txKey contextKey = iota
)

func New(db *dbx.DB) *DB {
	return &DB{db}
}

func (db *DB) DB() *dbx.DB {
	return db.db
}

func (db *DB) With(ctx context.Context) dbx.Builder {
	if tx, ok := ctx.Value(txKey).(*dbx.Tx); ok {
		return tx
	}
	return db.db.WithContext(ctx)
}

func (db *DB) Transactional(ctx context.Context, f func(ctx context.Context) error) error {
	return db.db.TransactionalContext(ctx, nil, func(tx *dbx.Tx) error {
		return f(context.WithValue(ctx, txKey, tx))
	})
}

func (db *DB) TransactionHandler() bunrouter.MiddlewareFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			return db.db.TransactionalContext(req.Context(), nil, func(tx *dbx.Tx) error {
				ctx := context.WithValue(req.Context(), txKey, tx)
				req = req.WithContext(ctx)
				return next(w, req)
			})
		}
	}
}
