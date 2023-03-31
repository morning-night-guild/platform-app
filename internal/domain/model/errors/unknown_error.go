package errors

import (
	"errors"
)

// UnknownError 未知のエラー.
type UnknownError struct {
	msg string
}

// NewUnknownError Unknownエラーのファクトリー関数.
func NewUnknownError(
	msg string,
) UnknownError {
	return UnknownError{
		msg: msg,
	}
}

// Error エラーメソッド.
func (ue UnknownError) Error() string {
	return ue.msg
}

// AsUnknownError UnknownError型に変換できるかどうかを判定する.
func AsUnknownError(err error) bool {
	var target UnknownError

	return errors.As(err, &target)
}
