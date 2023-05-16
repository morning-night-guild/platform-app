package auth

import (
	"net/mail"

	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

type Email string

func NewEmail(value string) (Email, error) {
	em := Email(value)

	if err := em.validate(); err != nil {
		return Email(""), err
	}

	return em, nil
}

func (em Email) String() string {
	return string(em)
}

func (em Email) validate() error {
	if _, err := mail.ParseAddress(em.String()); err != nil {
		log.Log().Warn("failed to parse email address", log.ErrorField(err))

		return errors.NewValidationError("invalid email address", err)
	}

	return nil
}
