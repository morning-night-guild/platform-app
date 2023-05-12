package model_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	type args struct {
		userID user.ID
	}

	tests := []struct {
		name string
		args args
		want model.User
	}{
		{
			name: "ユーザーを生成できる",
			args: args{
				userID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
			},
			want: model.User{
				UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := model.NewUser(tt.args.userID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
