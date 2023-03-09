package value

import (
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

// Index インデックス.
type Index int

// NewIndex インデックスファクトリー関数.
func NewIndex(value int) (Index, error) {
	if value < 0 {
		return Index(-1), errors.NewValidationError("index must be positive")
	}

	return Index(value), nil
}

// Int int型を取得する.
func (i Index) Int() int {
	return int(i)
}
