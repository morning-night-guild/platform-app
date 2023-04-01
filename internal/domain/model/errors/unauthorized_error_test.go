package errors_test

import (
	"fmt"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

func TestAsUnauthorizedError(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "UnauthorizedError型の場合はtrueを返す",
			args: args{
				err: errors.NewUnauthorizedError("test"),
			},
			want: true,
		},
		{
			name: "UnauthorizedError型の場合はtrueを返す",
			args: args{
				err: errors.NewUnauthorizedError("test", fmt.Errorf("test")),
			},
			want: true,
		},
		{
			name: "UnauthorizedError型ではない場合はfalseを返す",
			args: args{
				err: fmt.Errorf("test"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := errors.AsUnauthorizedError(tt.args.err); got != tt.want {
				t.Errorf("AsUnauthorizedError() = %v, want %v", got, tt.want)
			}
		})
	}
}
