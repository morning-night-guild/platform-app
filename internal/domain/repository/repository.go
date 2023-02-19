package repository

import (
	"fmt"

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
