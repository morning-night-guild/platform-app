package errors

import (
	"errors"
	"fmt"
)

// NotFoundError リソースが存在しないときに発生するエラー.
type NotFoundError struct {
	msg string
	err error
}

// NewNotFoundError NotFoundエラーのファクトリー関数.
func NewNotFoundError(
	msg string,
	errs ...error,
) NotFoundError {
	if len(errs) == 0 {
		return NotFoundError{
			msg: msg,
		}
	}

	return NotFoundError{
		msg: msg,
		err: errors.Join(errs...),
	}
}

// Error エラーメソッド.
func (err NotFoundError) Error() string {
	if err.err != nil {
		return fmt.Errorf("%s: %w", err.msg, err.err).Error()
	}

	return err.msg
}

// Unwrap アンラップ.
func (err NotFoundError) Unwrap() error {
	return err.err
}

// AsNotFoundError NotFoundError型に変換できるかどうかを判定する.
func AsNotFoundError(err error) bool {
	var target NotFoundError

	return errors.As(err, &target)
}
