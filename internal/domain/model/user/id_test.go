package user_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

func TestNewID(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    user.ID
		wantErr bool
	}{
		{
			name: "ユーザーIDが作成できる",
			args: args{
				value: "01234567-0123-0123-0123-0123456789ab",
			},
			want:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
			wantErr: false,
		},
		{
			name: "ユーザーIDが作成できない",
			args: args{
				value: "id",
			},
			want:    user.ID{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := user.NewID(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewID() = %v, want %v", got, tt.want)
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
			ctx := user.SetUIDCtx(tt.args.ctx, tt.args.uid)
			got := user.GetUIDCtx(ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got user id = %v, want %v", got, tt.want)
			}
		})
	}
}
