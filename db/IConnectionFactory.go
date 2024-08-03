package db

import (
	"context"
	"database/sql"
)

type IConnectionFactory interface {
	GetConnection(ctx context.Context) (*sql.Conn, error)
	GetTransaction(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
