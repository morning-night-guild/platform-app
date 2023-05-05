package external_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/adapter/external"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

func TestExternalHandleError(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		err error
	}

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(*testing.T, error)
	}{
		{
			name: "invalid argument",
			args: args{
				ctx: context.Background(),
				err: connect.NewError(connect.CodeInvalidArgument, nil),
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				if !errors.AsValidationError(err) {
					t.Errorf("External.HandleError() error = %v", err)
				}
			},
		},
		{
			name: "unauthorized",
			args: args{
				ctx: context.Background(),
				err: connect.NewError(connect.CodeUnauthenticated, nil),
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				if !errors.AsUnauthorizedError(err) {
					t.Errorf("External.HandleError() error = %v", err)
				}
			},
		},
		{
			name: "not found",
			args: args{
				ctx: context.Background(),
				err: connect.NewError(connect.CodeNotFound, nil),
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				if !errors.AsNotFoundError(err) {
					t.Errorf("External.HandleError() error = %v", err)
				}
			},
		},
		{
			name: "internal server",
			args: args{
				ctx: context.Background(),
				err: connect.NewError(connect.CodeInternal, nil),
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				if !errors.AsUnknownError(err) {
					t.Errorf("External.HandleError() error = %v", err)
				}
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				err: fmt.Errorf("error"),
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				if !errors.AsUnknownError(err) {
					t.Errorf("External.HandleError() error = %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ext := external.New()
			err := ext.HandleError(tt.args.ctx, tt.args.err)
			tt.checkErr(t, err)
		})
	}
}
