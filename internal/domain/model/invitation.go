package model

import (
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/notice"
)

type Invitation struct {
	Code  auth.InvitationCode
	Email auth.Email
}

func NewInvitation(
	code auth.InvitationCode,
	email auth.Email,
) (Invitation, error) {
	return Invitation{
		Code:  code,
		Email: email,
	}, nil
}

func GenerateInvitation(
	email auth.Email,
) Invitation {
	return Invitation{
		Code:  auth.GenerateInvitationCode(),
		Email: email,
	}
}

func (inv Invitation) Subject() notice.Subject {
	return notice.GenerateInvitationSubject()
}

func (inv Invitation) Message() notice.Message {
	return notice.GenerateInvitationMessage(inv.Code)
}
