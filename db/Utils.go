package db

import (
	"context"
	"database/sql"

	"github.com/uoul/go-common/async"
)

type ResultMapper[T any] func() ([]any, *T)
type EffectedRows int

func ExecStatement(connectionFactory IConnectionFactory, sql string, args ...any) chan async.ActionResult[EffectedRows] {
	result := make(chan async.ActionResult[EffectedRows])
	go func() {
		ctx := context.Background()
		conn, err := connectionFactory.GetConnection(ctx)
		if err != nil {
			result <- async.ActionResult[EffectedRows]{Result: 0, Error: err}
			return
		}
		defer conn.Close()
		r, err := conn.ExecContext(ctx, sql, args...)
		if err != nil {
			result <- async.ActionResult[EffectedRows]{Result: 0, Error: err}
			return
		}
		rowsEffected, err := r.RowsAffected()
		result <- async.ActionResult[EffectedRows]{Result: EffectedRows(rowsEffected), Error: err}
	}()
	return result
}

func ExecStatementTx(tx *sql.Tx, sql string, args ...any) chan async.ActionResult[EffectedRows] {
	result := make(chan async.ActionResult[EffectedRows])
	go func() {
		r, err := tx.Exec(sql, args...)
		if err != nil {
			result <- async.ActionResult[EffectedRows]{Result: 0, Error: err}
			return
		}
		rowsEffected, err := r.RowsAffected()
		result <- async.ActionResult[EffectedRows]{Result: EffectedRows(rowsEffected), Error: err}
	}()
	return result
}

func QueryStatement[T any](connectionFactory IConnectionFactory, resultMapper ResultMapper[T], sql string, args ...any) chan async.ActionResult[[]T] {
	result := make(chan async.ActionResult[[]T])
	go func(r chan async.ActionResult[[]T]) {
		ctx := context.Background()
		conn, err := connectionFactory.GetConnection(ctx)
		if err != nil {
			result <- async.ActionResult[[]T]{Result: []T{}, Error: err}
			return
		}
		defer conn.Close()
		rows, err := conn.QueryContext(ctx, sql, args...)
		if err != nil {
			result <- async.ActionResult[[]T]{Result: []T{}, Error: err}
			return
		}
		defer rows.Close()
		resultSet := []T{}
		for rows.Next() {
			fields, entry := resultMapper()
			err = rows.Scan(fields...)
			if err != nil {
				result <- async.ActionResult[[]T]{Result: []T{}, Error: err}
				return
			}
			resultSet = append(resultSet, *entry)
		}
		result <- async.ActionResult[[]T]{Result: resultSet, Error: nil}
	}(result)
	return result
}

func QueryStatementTx[T any](tx *sql.Tx, resultMapper ResultMapper[T], sql string, args ...any) chan async.ActionResult[[]T] {
	result := make(chan async.ActionResult[[]T])
	go func(r chan async.ActionResult[[]T]) {
		rows, err := tx.Query(sql, args...)
		if err != nil {
			result <- async.ActionResult[[]T]{Result: []T{}, Error: err}
			return
		}
		defer rows.Close()
		resultSet := []T{}
		for rows.Next() {
			fields, entry := resultMapper()
			err = rows.Scan(fields...)
			if err != nil {
				result <- async.ActionResult[[]T]{Result: []T{}, Error: err}
				return
			}
			resultSet = append(resultSet, *entry)
		}
		result <- async.ActionResult[[]T]{Result: resultSet, Error: nil}
	}(result)
	return result
}

func QuerySingle[T any](connectionFactory IConnectionFactory, resultMapper ResultMapper[T], sql string, args ...any) chan async.ActionResult[T] {
	result := make(chan async.ActionResult[T])
	go func() {
		ctx := context.Background()
		conn, err := connectionFactory.GetConnection(ctx)
		if err != nil {
			result <- async.ActionResult[T]{Result: *new(T), Error: err}
			return
		}
		defer conn.Close()
		row := conn.QueryRowContext(ctx, sql, args...)
		fields, entry := resultMapper()
		err = row.Scan(fields...)
		if err != nil {
			result <- async.ActionResult[T]{Result: *new(T), Error: err}
			return
		}
		result <- async.ActionResult[T]{Result: *entry, Error: nil}
	}()
	return result
}

func QuerySingleTx[T any](tx *sql.Tx, resultMapper ResultMapper[T], sql string, args ...any) chan async.ActionResult[T] {
	result := make(chan async.ActionResult[T])
	go func() {
		row := tx.QueryRow(sql, args...)
		fields, entry := resultMapper()
		err := row.Scan(fields...)
		if err != nil {
			result <- async.ActionResult[T]{Result: *new(T), Error: err}
			return
		}
		result <- async.ActionResult[T]{Result: *entry, Error: nil}
	}()
	return result
}

type TransactionScopeFunction[T any] func(ctx context.Context, tx *sql.Tx) (T, error)

func ExecInTransactionContext[T any](ctx context.Context, cf IConnectionFactory, tsf TransactionScopeFunction[T]) (T, error) {
	tx, err := cf.GetTransaction(ctx, nil)
	if err != nil {
		return *new(T), err
	}
	defer tx.Rollback()
	result, err := tsf(ctx, tx)
	if err != nil {
		return result, err
	}
	err = tx.Commit()
	if err != nil {
		return result, err
	}
	return result, nil
}

func EffectedRowsMapper() ([]any, *EffectedRows) {
	v := *new(EffectedRows)
	return []any{&v}, &v
}
