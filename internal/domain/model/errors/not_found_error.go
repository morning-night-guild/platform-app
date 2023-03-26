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
	err error,
) NotFoundError {
	return NotFoundError{
		msg: msg,
		err: err,
	}
}

// Error エラーメソッド.
func (nfe NotFoundError) Error() string {
	return fmt.Errorf("%s: %w", nfe.msg, nfe.err).Error()
}

// AsNotFoundError NotFoundError型に変換できるかどうかを判定する.
func AsNotFoundError(err error) bool {
	var target NotFoundError

	return errors.As(err, &target)
}
