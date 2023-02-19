package controller

import (
	"encoding/base64"
	"strconv"

	"github.com/morning-night-guild/platform-app/internal/domain/repository"
)

// ネクストトークン.
type NextToken string

// NewNextToken ネクストトークンを生成する.
func NewNextToken(value string) NextToken {
	return NextToken(value)
}

// CreateNextTokenFromIndex インデックスからネクストトークンを作成する.
func CreateNextTokenFromIndex(index repository.Index) NextToken {
	return NextToken(base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(index.Int()))))
}

// String ネクストトークンの文字列を提供する.
func (t NextToken) String() string {
	return string(t)
}

// ToIndex ネクストトークンからインデックスを作成する.
func (t NextToken) ToIndex() repository.Index {
	dec, err := base64.StdEncoding.DecodeString(t.String())
	if err != nil {
		return repository.Index(0)
	}

	index, err := strconv.Atoi(string(dec))
	if err != nil {
		return repository.Index(0)
	}

	return repository.Index(index)
}

// CreateNextToken ネクストトークンを作成する.
func (t NextToken) CreateNextToken(size repository.Size) NextToken {
	return NextToken(base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(t.ToIndex().Int() + size.Int()))))
}
