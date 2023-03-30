package interactor_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestCoreUserUpdateExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepository repository.User
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
				userRepository: &mock.UserRepository{
					T: t,
					User: model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					},
					SaveAssert: func(t *testing.T, item model.User) {
						t.Helper()

						if _, err := uuid.Parse(item.UserID.String()); err != nil {
							t.Errorf("UserRepository.Save() got = %v, err %v", item, err)
						}
					},
					FindAssert: func(t *testing.T, id user.ID) {
						t.Helper()

						want := user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab"))

						if !reflect.DeepEqual(id, want) {
							t.Errorf("UserRepository.Find() id = %v, want %v", id, want)
						}
					},
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
				userRepository: &mock.UserRepository{
					T: t,
					User: model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					},
					FindErr: fmt.Errorf("test"),
					FindAssert: func(t *testing.T, id user.ID) {
						t.Helper()
					},
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
				userRepository: &mock.UserRepository{
					T: t,
					User: model.User{
						UserID: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					},
					FindAssert: func(t *testing.T, id user.ID) {
						t.Helper()
					},
					SaveErr: fmt.Errorf("test"),
					SaveAssert: func(t *testing.T, item model.User) {
						t.Helper()
					},
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
				tt.fields.userRepository,
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
