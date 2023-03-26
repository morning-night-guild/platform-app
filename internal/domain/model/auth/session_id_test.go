package auth_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestNewSessionID(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    auth.SessionID
		wantErr bool
	}{
		{
			name: "セッションIDが作成できる",
			args: args{
				value: "01234567-0123-0123-0123-0123456789ab",
			},
			want:    auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
			wantErr: false,
		},
		{
			name: "セッションIDが作成できない",
			args: args{
				value: "id",
			},
			want:    auth.SessionID{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := auth.NewSessionID(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSessionID() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSessionID() = %v, want %v", got, tt.want)
			}
		})
	}
}
