package errors

import (
	"errors"
	"fmt"
)

// URLError URLに関するエラー.
type URLError struct {
	msg string
	err error
}

// NewURLError URLエラーのファクトリー関数.
func NewURLError(
	msg string,
	errs ...error,
) URLError {
	if len(errs) == 0 {
		return URLError{
			msg: msg,
		}
	}

	return URLError{
		msg: msg,
		err: errors.Join(errs...),
	}
}

// Error エラーメソッド.
func (err URLError) Error() string {
	if err.err != nil {
		return fmt.Errorf("%s: %w", err.msg, err.err).Error()
	}

	return err.msg
}

// Unwrap アンラップ.
func (err URLError) Unwrap() error {
	return err.err
}

// AsURLError URLError型に変換できるかどうかを判定する.
func AsURLError(err error) bool {
	var target URLError

	return errors.As(err, &target)
}
