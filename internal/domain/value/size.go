package value

import (
	"fmt"

	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

// Size サイズ.
type Size int

const maxSize = 100

// NewSize サイズファクトリー関数.
func NewSize(value int) (Size, error) {
	if value <= 0 {
		return Size(-1), errors.NewValidationError("size must be positive")
	}

	if value > maxSize {
		msg := fmt.Sprintf("size must be or less %d", maxSize)

		return Size(-1), errors.NewValidationError(msg)
	}

	return Size(value), nil
}

// Int int型を取得する.
func (s Size) Int() int {
	return int(s)
}
