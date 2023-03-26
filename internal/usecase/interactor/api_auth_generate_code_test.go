package interactor_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestAPIAuthGenerateCodeExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		secret    auth.Secret
		codeCache cache.Cache[model.Code]
	}

	type args struct {
		ctx   context.Context
		input port.APIAuthGenerateCodeInput
	}

	now := time.Now()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.APIAuthGenerateCodeOutput
		wantErr bool
	}{
		{
			name: "コードが生成できる",
			fields: fields{
				secret: auth.Secret("secret"),
				codeCache: &mock.Cache[model.Code]{
					T: t,
					SetAssert: func(t *testing.T, key string, value model.Code, ttl time.Duration) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIAuthGenerateCodeInput{
					SessionToken: auth.GenerateSessionToken(
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.Secret("secret"),
					),
				},
			},
			want: port.APIAuthGenerateCodeOutput{
				Code: model.Code{
					CodeID:    auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					IssuedAt:  now,
					ExpiresAt: now.Add(model.DefaultCodeExpiresIn),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aas := interactor.NewAPIAuthGenerateCode(
				tt.fields.secret,
				tt.fields.codeCache,
			)
			got, err := aas.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuthGenerateCode.Execute() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got.Code.SessionID, tt.want.Code.SessionID) {
				t.Errorf("APIAuthGenerateCode.Execute() got Code.SessionID = %v, want %v", got, tt.want)
			}
		})
	}
}
