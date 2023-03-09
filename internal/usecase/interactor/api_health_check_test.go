package interactor_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestAPIHealthCheckExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		healthRPC rpc.Health
	}

	type args struct {
		ctx   context.Context
		input port.APIHealthCheckInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.APIHealthCheckOutput
		wantErr bool
	}{
		{
			name: "ヘルスチェックが成功する",
			fields: fields{
				healthRPC: &mock.RPCHealth{
					T:   t,
					Err: nil,
				},
			},
			args: args{
				ctx:   context.Background(),
				input: port.APIHealthCheckInput{},
			},
			want:    port.APIHealthCheckOutput{},
			wantErr: false,
		},
		{
			name: "ヘルスチェックが失敗する",
			fields: fields{
				healthRPC: &mock.RPCHealth{
					T:   t,
					Err: errors.New("error"),
				},
			},
			args: args{
				ctx:   context.Background(),
				input: port.APIHealthCheckInput{},
			},
			want:    port.APIHealthCheckOutput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ahc := interactor.NewAPIHealthCheck(tt.fields.healthRPC)
			got, err := ahc.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIHealthCheck.Execute() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIHealthCheck.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
