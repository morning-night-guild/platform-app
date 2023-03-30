package errors

import (
	"errors"
	"fmt"
)

// UnknownError 未知のエラー.
type UnknownError struct {
	msg string
	err error
}

// NewUnknownError Unknownエラーのファクトリー関数.
func NewUnknownError(
	msg string,
	err error,
) UnknownError {
	return UnknownError{
		msg: msg,
		err: err,
	}
}

// Error エラーメソッド.
func (nfe UnknownError) Error() string {
	return fmt.Errorf("%s: %w", nfe.msg, nfe.err).Error()
}

// AsUnknownError UnknownError型に変換できるかどうかを判定する.
func AsUnknownError(err error) bool {
	var target UnknownError

	return errors.As(err, &target)
}
