package auth_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestNewEMail(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    auth.EMail
		wantErr bool
	}{
		{
			name: "メールアドレスが作成できる",
			args: args{
				value: "test@example.com",
			},
			want:    auth.EMail("test@example.com"),
			wantErr: false,
		},
		{
			name: "メールアドレスが作成できない",
			args: args{
				value: "email",
			},
			want:    auth.EMail(""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := auth.NewEMail(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEMail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewEMail() = %v, want %v", got, tt.want)
			}
		})
	}
}
