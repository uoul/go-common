package db

import "fmt"

type DbConnectionError struct {
	msg string
	err error
}

func NewDbConnectionError(msg string, err error) *DbConnectionError {
	return &DbConnectionError{
		msg: msg,
		err: err,
	}
}

func (e *DbConnectionError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}
