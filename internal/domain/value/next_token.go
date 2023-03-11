package value

import (
	"encoding/base64"
	"strconv"
)

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
