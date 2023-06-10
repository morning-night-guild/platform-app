package notice_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/notice"
)

func TestGenerateInvitationSubject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want notice.Subject
	}{
		{
			name: "招待メールの件名が生成できる",
			want: notice.Subject("Welcome to Morning Night Guild Platform!"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := notice.GenerateInvitationSubject(); got != tt.want {
				t.Errorf("GenerateInvitationSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}
