package auth_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestNewEmail(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    auth.Email
		wantErr bool
	}{
		{
			name: "メールアドレスが作成できる",
			args: args{
				value: "test@example.com",
			},
			want:    auth.Email("test@example.com"),
			wantErr: false,
		},
		{
			name: "メールアドレスが作成できない",
			args: args{
				value: "email",
			},
			want:    auth.Email(""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := auth.NewEmail(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
