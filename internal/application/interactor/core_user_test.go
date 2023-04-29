package interactor_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/application/interactor"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
)

const uid = "01234567-0123-0123-0123-0123456789ab"

func TestCoreUserCreate(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepository func(t *testing.T) repository.User
	}

	type args struct {
		ctx   context.Context
		input usecase.CoreUserCreateInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.CoreUserCreateOutput
		wantErr bool
	}{
		{
			name: "ユーザーが作成できる",
			fields: fields{
				userRepository: func(t *testing.T) repository.User {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockUser(ctrl)
					mock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
					return mock
				},
			},
			args: args{
				ctx:   context.Background(),
				input: usecase.CoreUserCreateInput{},
			},
			want:    usecase.CoreUserCreateOutput{},
			wantErr: false,
		},
		{
			name: "UserRepository.Save()でエラーが発生した場合はエラーを返す",
			fields: fields{
				userRepository: func(t *testing.T) repository.User {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := repository.NewMockUser(ctrl)
					mock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				ctx:   context.Background(),
				input: usecase.CoreUserCreateInput{},
			},
			want:    usecase.CoreUserCreateOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewCoreUser(tt.fields.userRepository(t))
			got, err := itr.Create(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoreUser.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, err := uuid.Parse(got.User.UserID.String()); err != nil {
				t.Errorf("CoreUser.Create() got User.UserID = %v, err %v", got, err)
			}
		})
	}
}

func TestCoreUserUpdate(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepository func(t *testing.T) repository.User
	}

	type args struct {
		ctx   context.Context
		input usecase.CoreUserUpdateInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.CoreUserUpdateOutput
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
						UserID: user.ID(uuid.MustParse(uid)),
					}, nil)
					mock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreUserUpdateInput{
					UserID: user.ID(uuid.MustParse(uid)),
				},
			},
			want: usecase.CoreUserUpdateOutput{
				User: model.User{
					UserID: user.ID(uuid.MustParse(uid)),
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
					mock.EXPECT().Find(gomock.Any(), user.ID(uuid.MustParse(uid))).Return(model.User{}, fmt.Errorf("test"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreUserUpdateInput{
					UserID: user.ID(uuid.MustParse(uid)),
				},
			},
			want: usecase.CoreUserUpdateOutput{
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
					mock.EXPECT().Find(gomock.Any(), user.ID(uuid.MustParse(uid))).Return(model.User{
						UserID: user.ID(uuid.MustParse(uid)),
					}, nil)
					mock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(fmt.Errorf("test"))
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				input: usecase.CoreUserUpdateInput{
					UserID: user.ID(uuid.MustParse(uid)),
				},
			},
			want: usecase.CoreUserUpdateOutput{
				User: model.User{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewCoreUser(tt.fields.userRepository(t))
			got, err := itr.Update(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoreUser.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CoreUser.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
