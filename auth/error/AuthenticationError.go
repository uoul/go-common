package error

import "fmt"

type AuthenticationError struct {
	msg string
}

func NewAuthenticationError(format string, a ...any) error {
	return &AuthenticationError{
		msg: fmt.Sprintf(format, a...),
	}
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("AuthenticationError | %s", e.msg)
}
