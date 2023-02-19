package api_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/adapter/api"
)

func TestAPIHandleConnectError(t *testing.T) {
	t.Parallel()

	type fields struct {
		key     string
		connect *api.Connect
	}

	type args struct {
		ctx context.Context
		err error
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "400",
			fields: fields{
				key:     "",
				connect: nil,
			},
			args: args{
				ctx: context.Background(),
				err: connect.NewError(connect.CodeInvalidArgument, errors.New("invalid argument")),
			},
			want: http.StatusBadRequest,
		},
		{
			name: "500",
			fields: fields{
				key:     "",
				connect: nil,
			},
			args: args{
				ctx: context.Background(),
				err: errors.New("unknown error"),
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			api := api.New(tt.fields.key, tt.fields.connect)
			if got := api.HandleConnectError(tt.args.ctx, tt.args.err); got != tt.want {
				t.Errorf("API.HandleConnectError() = %v, want %v", got, tt.want)
			}
		})
	}
}
