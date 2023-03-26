package auth_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestNewPassword(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    auth.Password
		wantErr bool
	}{
		{
			name: "パスワードが作成できる",
			args: args{
				value: "password",
			},
			want:    auth.Password("password"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := auth.NewPassword(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("NewPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
