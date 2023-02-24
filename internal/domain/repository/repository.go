package repository

import (
	"encoding/base64"
	"fmt"
	"strconv"

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

// ネクストトークン.
type NextToken string

// NewNextToken ネクストトークンを生成する.
func NewNextToken(value string) NextToken {
	return NextToken(value)
}

// CreateNextTokenFromIndex インデックスからネクストトークンを作成する.
func CreateNextTokenFromIndex(index Index) NextToken {
	return NextToken(base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(index.Int()))))
}

// String ネクストトークンの文字列を提供する.
func (t NextToken) String() string {
	return string(t)
}

// ToIndex ネクストトークンからインデックスを作成する.
func (t NextToken) ToIndex() Index {
	dec, err := base64.StdEncoding.DecodeString(t.String())
	if err != nil {
		return Index(0)
	}

	index, err := strconv.Atoi(string(dec))
	if err != nil {
		return Index(0)
	}

	return Index(index)
}

// CreateNextToken ネクストトークンを作成する.
func (t NextToken) CreateNextToken(size Size) NextToken {
	return NextToken(base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(t.ToIndex().Int() + size.Int()))))
}
