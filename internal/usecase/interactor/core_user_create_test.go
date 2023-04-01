package interactor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestCoreUserCreateExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepository func(t *testing.T) repository.User
	}

	type args struct {
		ctx   context.Context
		input port.CoreUserCreateInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.CoreUserCreateOutput
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
				input: port.CoreUserCreateInput{},
			},
			want:    port.CoreUserCreateOutput{},
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
				input: port.CoreUserCreateInput{},
			},
			want:    port.CoreUserCreateOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cuc := interactor.NewCoreUserCreate(tt.fields.userRepository(t))
			got, err := cuc.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoreUserCreate.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, err := uuid.Parse(got.User.UserID.String()); err != nil {
				t.Errorf("CoreUserCreate.Execute() got User.UserID = %v, err %v", got, err)
			}
		})
	}
}
