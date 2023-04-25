package auth

import (
	"fmt"
	"time"

	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

const (
	DefaultExpiresInSecond = 3600 // 1 hour (60s * 60m)
	DefaultExpiresIn       = ExpiresIn(DefaultExpiresInSecond)
)

type ExpiresIn int

func NewExpiresIn(value int) (ExpiresIn, error) {
	ei := ExpiresIn(value)

	if err := ei.validate(); err != nil {
		return ExpiresIn(-1), err
	}

	return ei, nil
}

func (ei ExpiresIn) Int() int {
	return int(ei)
}

func (ei ExpiresIn) Duration() time.Duration {
	return time.Duration(ei.Int()) * time.Second
}

func (ei ExpiresIn) validate() error {
	if ei < 0 {
		msg := fmt.Sprintf("ExpiresIn must be greater than or equal to 0: %d", ei)

		return errors.NewValidationError(msg)
	}

	if ei > DefaultExpiresInSecond {
		msg := fmt.Sprintf("ExpiresIn must be less than or equal to %d: %d", DefaultExpiresInSecond, ei)

		return errors.NewValidationError(msg)
	}

	return nil
}
