package errors

import "errors"

// ValidationError 値オブジェクト生成時に発生するバリデーションエラー.
type ValidationError struct {
	msg string
}

// NewValidationError バリデーションエラーのファクトリー関数.
func NewValidationError(msg string) ValidationError {
	return ValidationError{
		msg: msg,
	}
}

// Error エラーメソッド.
func (ve ValidationError) Error() string {
	return ve.msg
}

// AsValidationError ValidationError型に変換できるかどうかを判定する.
func AsValidationError(err error) bool {
	var target ValidationError

	return errors.As(err, &target)
}
