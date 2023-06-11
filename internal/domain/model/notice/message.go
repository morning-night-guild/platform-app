package notice

import (
	"fmt"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

type Message string

func (msg Message) String() string {
	return string(msg)
}

const invitationMessage = `Welcome to Morning Night Guild Platform!

Invitation Code
====================
%s
====================

Please enter the invitation code on the registration screen.
`

func GenerateInvitationMessage(code auth.InvitationCode) Message {
	msg := fmt.Sprintf(invitationMessage, code.String())

	return Message(msg)
}
