package model_test

import (
	"context"
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

func TestSetUIDCtxAndGetUIDCtx(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		uid user.ID
	}

	tests := []struct {
		name string
		args args
		want user.ID
	}{
		{
			name: "contextにUserIDを設定できる",
			args: args{
				ctx: context.Background(),
				uid: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
			},
			want: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
		},
		{
			name: "contextにUserIDが設定されていない",
			args: args{
				ctx: context.Background(),
				uid: user.GenerateZeroID(),
			},
			want: user.GenerateZeroID(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := model.SetUIDCtx(tt.args.ctx, tt.args.uid)
			got := model.GetUIDCtx(ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got user id = %v, want %v", got, tt.want)
			}
		})
	}
}
