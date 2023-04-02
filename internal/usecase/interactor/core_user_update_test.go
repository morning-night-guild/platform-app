package interactor_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestCoreUserUpdateExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepository func(t *testing.T) repository.User
	}

	type args struct {
		ctx   context.Context
		input port.CoreUserUpdateInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.CoreUserUpdateOutput
		wantErr bool
	}{
		{
			name: "ユーザーが更新できる",
			fields: fields{
				userRepository: func(t *testing.T) repository.User {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockUser(ctrl)
					mock.EXPECT().Find(gomock.Any(), gomock.Any()).Return(model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					}, nil)
					mock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.CoreUserUpdateInput{
					UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				},
			},
			want: port.CoreUserUpdateOutput{
				User: model.User{
					UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				},
			},
			wantErr: false,
		},
		{
			name: "UserRepository.Find()でエラーが発生してユーザーが更新できない",
			fields: fields{
				userRepository: func(t *testing.T) repository.User {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockUser(ctrl)
					mock.EXPECT().Find(gomock.Any(), user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab"))).Return(model.User{}, fmt.Errorf("test"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.CoreUserUpdateInput{
					UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				},
			},
			want: port.CoreUserUpdateOutput{
				User: model.User{},
			},
			wantErr: true,
		},
		{
			name: "UserRepository.Save()でエラーが発生してユーザーが更新できない",
			fields: fields{
				userRepository: func(t *testing.T) repository.User {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockUser(ctrl)
					mock.EXPECT().Find(gomock.Any(), user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab"))).Return(model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					}, nil)
					mock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(fmt.Errorf("test"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.CoreUserUpdateInput{
					UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				},
			},
			want: port.CoreUserUpdateOutput{
				User: model.User{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cuu := interactor.NewCoreUserUpdate(
				tt.fields.userRepository(t),
			)
			got, err := cuu.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoreUserUpdate.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CoreUserUpdate.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
