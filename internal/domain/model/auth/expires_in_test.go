package auth_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestNewExpiresIn(t *testing.T) {
	t.Parallel()

	type args struct {
		value int
	}

	tests := []struct {
		name    string
		args    args
		want    auth.ExpiresIn
		wantErr bool
	}{
		{
			name: "最小値で有効期限(0s)が作成できる",
			args: args{
				value: 0,
			},
			want:    auth.ExpiresIn(0),
			wantErr: false,
		},
		{
			name: "有効期限(3600s)が作成できる",
			args: args{
				value: 3600,
			},
			want:    auth.ExpiresIn(3600),
			wantErr: false,
		},
		{
			name: "最大値で有効期限(86400s)が作成できる",
			args: args{
				value: 86400,
			},
			want:    auth.ExpiresIn(86400),
			wantErr: false,
		},
		{
			name: "負の値で有効期限が作成できない",
			args: args{
				value: -1,
			},
			want:    auth.ExpiresIn(-1),
			wantErr: true,
		},
		{
			name: "最大値で有効期限(86400s)より大きい値で作成できない",
			args: args{
				value: 86401,
			},
			want:    auth.ExpiresIn(-1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := auth.NewExpiresIn(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExpiresIn() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("NewExpiresIn() = %v, want %v", got, tt.want)
			}
		})
	}
}
