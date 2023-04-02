package auth_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestNewSignature(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    auth.Signature
		wantErr bool
	}{
		{
			name: "署名が作成できる",
			args: args{
				value: "signature",
			},
			want:    auth.Signature("signature"),
			wantErr: false,
		},
		{
			name: "署名が作成できない",
			args: args{
				value: "",
			},
			want:    auth.Signature(""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := auth.NewSignature(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
