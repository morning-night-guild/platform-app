package errors

import (
	"errors"
)

// UnauthorizedError 認証エラー.
type UnauthorizedError struct {
	msg string
}

// NewUnauthorizedError 認証エラーのファクトリー関数.
func NewUnauthorizedError(
	msg string,
) UnauthorizedError {
	return UnauthorizedError{
		msg: msg,
	}
}

// Error エラーメソッド.
func (ue UnauthorizedError) Error() string {
	return ue.msg
}

// AsUnauthorizedError UnauthorizedError型に変換できるかどうかを判定する.
func AsUnauthorizedError(err error) bool {
	var target UnauthorizedError

	return errors.As(err, &target)
}
