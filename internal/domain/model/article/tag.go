package article

import (
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

// Tag 記事のタグ.
type Tag string

// NewTag タグを作成するファクトリー関数.
func NewTag(value string) (Tag, error) {
	tag := Tag(value)

	if err := tag.validate(); err != nil {
		return Tag(""), err
	}

	return tag, nil
}

// String 記事のタグを文字列として提供するメソッド.
func (t Tag) String() string {
	return string(t)
}

// validate 記事のタグを検証するメソッド.
func (t Tag) validate() error {
	if len(t) == 0 {
		return errors.NewValidationError("tag must not be empty")
	}

	return nil
}
