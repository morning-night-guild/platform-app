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

	type args struct {
		ctx context.Context
		err error
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "400",
			args: args{
				ctx: context.Background(),
				err: connect.NewError(connect.CodeInvalidArgument, errors.New("invalid argument")),
			},
			want: http.StatusBadRequest,
		},
		{
			name: "500",
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
			api := api.New("", nil, nil)
			if got := api.HandleConnectError(tt.args.ctx, tt.args.err); got != tt.want {
				t.Errorf("API.HandleConnectError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIPointerToString(t *testing.T) {
	t.Parallel()

	type args struct {
		s *string
	}

	test := "test"

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil",
			args: args{
				s: nil,
			},
			want: "",
		},
		{
			name: "not nil",
			args: args{
				s: &test,
			},
			want: "test",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			api := api.New("", nil, nil)
			if got := api.PointerToString(tt.args.s); got != tt.want {
				t.Errorf("API.PointerToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
