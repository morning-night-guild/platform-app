package article

import (
	"fmt"
	"strings"

	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

// Thumbnail サムネイル.
type Thumbnail string

// String サムネイルを文字列として提供するメソッド.
func (t Thumbnail) String() string {
	return string(t)
}

// NewThumbnail サムネイルを新規作成するファクトリー関数.
func NewThumbnail(t string) (Thumbnail, error) {
	thumbnail := Thumbnail(t)

	if err := thumbnail.validate(); err != nil {
		return Thumbnail(""), err
	}

	return thumbnail, nil
}

// validate サムネイルを検証するメソッド.
func (t Thumbnail) validate() error {
	// 空文字も扱うこととし
	// 空文字だったらバリデーションを実施しない
	if len(t.String()) == 0 {
		return nil
	}

	if !strings.HasPrefix(t.String(), "https://") {
		msg := fmt.Sprintf("thumbnail must be start with `https://`. value is %s", t.String())

		return errors.NewValidationError(msg)
	}

	return nil
}
