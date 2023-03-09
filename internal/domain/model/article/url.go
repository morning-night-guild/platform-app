package article

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

// URL 記事のURL.
type URL string

// String URLを文字列として提供するメソッド.
func (ur URL) String() string {
	return string(ur)
}

// NewURL URLを新規作成する関数.
func NewURL(value string) (URL, error) {
	ur := URL(value)

	if err := ur.validate(); err != nil {
		return URL(""), err
	}

	return ur, nil
}

// validate URLを検証するメソッド.
func (ur URL) validate() error {
	if _, err := url.Parse(ur.String()); err != nil {
		return errors.NewValidationError(err.Error())
	}

	if !strings.HasPrefix(ur.String(), "https://") {
		msg := fmt.Sprintf("url must be start with `https://`. value is %s", ur.String())

		return errors.NewValidationError(msg)
	}

	return nil
}
