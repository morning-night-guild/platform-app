package errors_test

import (
	"fmt"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

func TestAsURLError(t *testing.T) {
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
			name: "URLError型の場合はtrueを返す",
			args: args{
				err: errors.NewURLError("test"),
			},
			want: true,
		},
		{
			name: "URLError型ではない場合はfalseを返す",
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
			if got := errors.AsURLError(tt.args.err); got != tt.want {
				t.Errorf("AsURLError() = %v, want %v", got, tt.want)
			}
		})
	}
}
