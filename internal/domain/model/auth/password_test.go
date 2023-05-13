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
		{
			name: "パスワードが作成できない",
			args: args{
				value: "",
			},
			want:    auth.Password(""),
			wantErr: true,
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

func TestPasswordEqual(t *testing.T) {
	t.Parallel()

	type args struct {
		password auth.Password
	}

	tests := []struct {
		name string
		pw   auth.Password
		args args
		want bool
	}{
		{
			name: "パスワードが一致する",
			pw:   auth.Password("password"),
			args: args{
				password: auth.Password("password"),
			},
			want: true,
		},
		{
			name: "パスワードが一致しない",
			pw:   auth.Password("password"),
			args: args{
				password: auth.Password("Password"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.pw.Equal(tt.args.password); got != tt.want {
				t.Errorf("Password.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
