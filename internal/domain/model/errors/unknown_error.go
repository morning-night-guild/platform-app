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
	errs ...error,
) UnknownError {
	if len(errs) == 0 {
		return UnknownError{
			msg: msg,
		}
	}

	return UnknownError{
		msg: msg,
		err: errors.Join(errs...),
	}
}

// Error エラーメソッド.
func (err UnknownError) Error() string {
	if err.err != nil {
		return fmt.Errorf("%s: %w", err.msg, err.err).Error()
	}

	return err.msg
}

// Unwrap アンラップ.
func (err UnknownError) Unwrap() error {
	return err.err
}

// AsUnknownError UnknownError型に変換できるかどうかを判定する.
func AsUnknownError(err error) bool {
	var target UnknownError

	return errors.As(err, &target)
}
