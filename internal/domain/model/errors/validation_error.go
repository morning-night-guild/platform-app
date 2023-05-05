package errors

import (
	"errors"
	"fmt"
)

// ValidationError 値オブジェクト生成時に発生するバリデーションエラー.
type ValidationError struct {
	msg string
	err error
}

// NewValidationError バリデーションエラーのファクトリー関数.
func NewValidationError(
	msg string,
	errs ...error,
) ValidationError {
	if len(errs) == 0 {
		return ValidationError{
			msg: msg,
		}
	}

	return ValidationError{
		msg: msg,
		err: errors.Join(errs...),
	}
}

// Error エラーメソッド.
func (ve ValidationError) Error() string {
	if ve.err != nil {
		return fmt.Errorf("%s: %w", ve.msg, ve.err).Error()
	}

	return ve.msg
}

// AsValidationError ValidationError型に変換できるかどうかを判定する.
func AsValidationError(err error) bool {
	var target ValidationError

	return errors.As(err, &target)
}
