package notice_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/notice"
)

const wantInvitationCodeMessage = `Welcome to Morning Night Guild Platform!

Invitation Code
====================
test
====================

Please enter the invitation code on the registration screen.
`

func TestGenerateInvitationMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		code auth.InvitationCode
	}

	tests := []struct {
		name string
		args args
		want notice.Message
	}{
		{
			name: "招待メッセージが生成できる",
			args: args{
				code: auth.InvitationCode("test"),
			},
			want: notice.Message(wantInvitationCodeMessage),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := notice.GenerateInvitationMessage(tt.args.code); got != tt.want {
				t.Errorf("GenerateInvitationMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
