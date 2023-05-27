package article

import "github.com/morning-night-guild/platform-app/internal/domain/model/errors"

type Scope string

const (
	All Scope = "all"
	Own Scope = "own"
)

func NewScope(s string) (Scope, error) {
	if s != All.String() && s != Own.String() {
		return "", errors.NewValidationError("invalid scope")
	}

	return Scope(s), nil
}

func (sc Scope) String() string {
	return string(sc)
}
