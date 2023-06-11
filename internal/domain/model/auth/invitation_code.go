package auth

import (
	"strings"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

type InvitationCode string

func (ic InvitationCode) String() string {
	return string(ic)
}

func (ic InvitationCode) validate() error {
	if ic == "" {
		return errors.NewValidationError("invitation code is required")
	}

	return nil
}

func NewInvitationCode(value string) (InvitationCode, error) {
	code := InvitationCode(value)

	if err := code.validate(); err != nil {
		return InvitationCode(""), err
	}

	return code, nil
}

func GenerateInvitationCode() InvitationCode {
	code := strings.Split(uuid.NewString(), "-")[0]

	return InvitationCode(code)
}
