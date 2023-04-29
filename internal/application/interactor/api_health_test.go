package interactor_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/morning-night-guild/platform-app/internal/application/interactor"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
)

func TestAPIHealthCheck(t *testing.T) {
	t.Parallel()

	type fields struct {
		healthRPC func(t *testing.T) rpc.Health
	}

	type args struct {
		ctx   context.Context
		input usecase.APIHealthCheckInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.APIHealthCheckOutput
		wantErr bool
	}{
		{
			name: "ヘルスチェックが成功する",
			fields: fields{
				healthRPC: func(t *testing.T) rpc.Health {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockHealth(ctrl)
					mock.EXPECT().Check(gomock.Any()).Return(nil)
					return mock
				},
			},
			args: args{
				ctx:   context.Background(),
				input: usecase.APIHealthCheckInput{},
			},
			want:    usecase.APIHealthCheckOutput{},
			wantErr: false,
		},
		{
			name: "ヘルスチェックが失敗する",
			fields: fields{
				healthRPC: func(t *testing.T) rpc.Health {
					t.Helper()
					ctrl := gomock.NewController(t)
					mock := rpc.NewMockHealth(ctrl)
					mock.EXPECT().Check(gomock.Any()).Return(fmt.Errorf("error"))
					return mock
				},
			},
			args: args{
				ctx:   context.Background(),
				input: usecase.APIHealthCheckInput{},
			},
			want:    usecase.APIHealthCheckOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			itr := interactor.NewAPIHealth(tt.fields.healthRPC(t))
			got, err := itr.Check(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIHealth.Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIHealth.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
