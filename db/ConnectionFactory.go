package db

import (
	"context"
	"database/sql"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
)

type ConnectionFactory struct {
	mux           sync.Mutex
	connectionStr string
	driver        string
	db            *sql.DB
}

// GetTransaction implements IConnectionFactory.
func (f *ConnectionFactory) GetTransaction(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	db, err := f.getDb(ctx)
	if err != nil {
		return nil, err
	}
	return db.BeginTx(ctx, opts)
}

// GetConnection implements IConnectionFactory.
func (f *ConnectionFactory) GetConnection(ctx context.Context) (*sql.Conn, error) {
	db, err := f.getDb(ctx)
	if err != nil {
		return nil, err
	}
	return db.Conn(ctx)
}

func (f *ConnectionFactory) getDb(ctx context.Context) (*sql.DB, error) {
	if f.db == nil {
		f.mux.Lock()
		defer f.mux.Unlock()
		if f.db == nil {
			db, err := sql.Open(f.driver, f.connectionStr)
			if err != nil {
				return nil, NewDbConnectionError("failed to open database", err)
			}
			f.db = db
		}
	} else {
		if err := f.db.PingContext(ctx); err != nil {
			f.db = nil
			return nil, err
		}
	}
	return f.db, nil
}

func NewConnectionFactory(connectionString, driver string) IConnectionFactory {
	return &ConnectionFactory{
		connectionStr: connectionString,
		driver:        driver,
	}
}
