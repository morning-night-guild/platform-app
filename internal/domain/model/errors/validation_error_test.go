package errors_test

import (
	"fmt"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

func TestAsValidationError(t *testing.T) {
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
			name: "ValidationError型の場合はtrueを返す",
			args: args{
				err: errors.NewValidationError("test"),
			},
			want: true,
		},
		{
			name: "ValidationError型の場合はtrueを返す",
			args: args{
				err: errors.NewValidationError("test", fmt.Errorf("test")),
			},
			want: true,
		},
		{
			name: "ValidationError型ではない場合はfalseを返す",
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
			if got := errors.AsValidationError(tt.args.err); got != tt.want {
				t.Errorf("AsValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}
