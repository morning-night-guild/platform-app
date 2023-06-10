package model_test

import (
	"reflect"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestGenerateInvitation(t *testing.T) {
	t.Parallel()

	type args struct {
		email auth.Email
	}

	tests := []struct {
		name string
		args args
		want model.Invitation
	}{
		{
			name: "招待情報が生成できる",
			args: args{
				email: auth.Email("test@example.com"),
			},
			want: model.Invitation{
				Code:  auth.InvitationCode(""), // NOTE: 本来はランダムな文字列が入る
				Email: auth.Email("test@example.com"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := model.GenerateInvitation(tt.args.email)
			if !reflect.DeepEqual(got.Email, tt.want.Email) {
				t.Errorf("GenerateInvitation() = %v, want %v", got, tt.want)
			}
			if len(got.Code) != 8 {
				t.Errorf("GenerateInvitation() = %v, want %v", got, tt.want)
			}
		})
	}
}
