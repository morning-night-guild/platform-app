package errors

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
func (e URLError) Error() string {
	return e.msg
}
