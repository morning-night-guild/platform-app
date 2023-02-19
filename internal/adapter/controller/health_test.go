package controller_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/adapter/controller"
	healthv1 "github.com/morning-night-guild/platform-app/pkg/connect/proto/health/v1"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *connect.Request[healthv1.CheckRequest]
	}

	tests := []struct {
		name    string
		h       controller.Health
		args    args
		want    *connect.Response[healthv1.CheckResponse]
		wantErr bool
	}{
		{
			name: "ヘルスチェックができる",
			args: args{
				ctx: context.Background(),
				req: &connect.Request[healthv1.CheckRequest]{},
			},
			want:    connect.NewResponse(&healthv1.CheckResponse{}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := controller.NewHealth()
			got, err := h.Check(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Health.Check() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Health.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
