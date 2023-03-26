package errors

import "errors"

// URLError URLに関するエラー.
type URLError struct {
	msg string
}

// NewURLError URLエラーのファクトリー関数.
func NewURLError(msg string) URLError {
	return URLError{
		msg: msg,
	}
}

// Error エラーメソッド.
func (ue URLError) Error() string {
	return ue.msg
}

// AsURLError URLError型に変換できるかどうかを判定する.
func AsURLError(err error) bool {
	var target URLError

	return errors.As(err, &target)
}
