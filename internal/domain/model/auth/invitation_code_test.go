package auth_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestGenerateInvitationCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "招待コードが生成できる",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := auth.GenerateInvitationCode(); len(got) != 8 {
				t.Errorf("GenerateInvitationCode() = %v", got)
			}
		})
	}
}
