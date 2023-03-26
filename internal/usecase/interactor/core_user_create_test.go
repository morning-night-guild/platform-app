package interactor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestCoreUserCreateExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepository repository.User
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
				userRepository: &mock.UserRepository{
					T: t,
					User: model.User{
						UserID: user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					},
					SaveAssert: func(t *testing.T, item model.User) {
						if _, err := uuid.Parse(item.UserID.String()); err != nil {
							t.Errorf("UserRepository.Save() got = %v, err %v", item, err)
						}
					},
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
			name: "UserReository.Save()でエラーが発生した場合はエラーを返す",
			fields: fields{
				userRepository: &mock.UserRepository{
					T:          t,
					User:       model.User{},
					Err:        fmt.Errorf("test"),
					SaveAssert: func(t *testing.T, item model.User) {},
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
			cuc := interactor.NewCoreUserCreate(tt.fields.userRepository)
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
