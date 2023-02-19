package errors

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
func (e ValidationError) Error() string {
	return e.msg
}
