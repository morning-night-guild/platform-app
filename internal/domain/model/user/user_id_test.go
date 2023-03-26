package user_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

func TestNewUserID(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    user.UserID
		wantErr bool
	}{
		{
			name: "ユーザーIDが作成できる",
			args: args{
				value: "01234567-0123-0123-0123-0123456789ab",
			},
			want:    user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
			wantErr: false,
		},
		{
			name: "ユーザーIDが作成できない",
			args: args{
				value: "id",
			},
			want:    user.UserID{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := user.NewUserID(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserID() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}
