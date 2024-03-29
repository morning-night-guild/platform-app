package auth

import "github.com/morning-night-guild/platform-app/internal/domain/model/errors"

type Password string

func NewPassword(value string) (Password, error) {
	pw := Password(value)

	if err := pw.validate(); err != nil {
		return Password(""), err
	}

	return pw, nil
}

func (pw Password) String() string {
	return string(pw)
}

func (pw Password) validate() error {
	if pw == "" {
		return errors.NewValidationError("password is required")
	}

	return nil
}

func (pw Password) Equal(password Password) bool {
	return pw.String() == password.String()
}
