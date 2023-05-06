package errors

import (
	"errors"
	"fmt"
)

// UnauthorizedError 認証エラー.
type UnauthorizedError struct {
	msg string
	err error
}

// NewUnauthorizedError 認証エラーのファクトリー関数.
func NewUnauthorizedError(
	msg string,
	errs ...error,
) UnauthorizedError {
	if len(errs) == 0 {
		return UnauthorizedError{
			msg: msg,
		}
	}

	return UnauthorizedError{
		msg: msg,
		err: errors.Join(errs...),
	}
}

// Error エラーメソッド.
func (err UnauthorizedError) Error() string {
	if err.err != nil {
		return fmt.Errorf("%s: %w", err.msg, err.err).Error()
	}

	return err.msg
}

// Unwrap アンラップ.
func (err UnauthorizedError) Unwrap() error {
	return err.err
}

// AsUnauthorizedError UnauthorizedError型に変換できるかどうかを判定する.
func AsUnauthorizedError(err error) bool {
	var target UnauthorizedError

	return errors.As(err, &target)
}
