package auth

import (
	"net/mail"

	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

type EMail string

func NewEMail(value string) (EMail, error) {
	em := EMail(value)

	if err := em.validate(); err != nil {
		return EMail(""), err
	}

	return em, nil
}

func (em EMail) String() string {
	return string(em)
}

func (em EMail) validate() error {
	if _, err := mail.ParseAddress(em.String()); err != nil {
		log.Log().Warn("failed to parse email address", log.ErrorField(err))

		return errors.NewValidationError("invalid email address", err)
	}

	return nil
}
