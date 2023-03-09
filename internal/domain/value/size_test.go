package value_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

func TestNewSize(t *testing.T) {
	t.Parallel()

	type args struct {
		value int
	}

	tests := []struct {
		name    string
		args    args
		want    value.Size
		wantErr bool
	}{
		{
			name: "1で作成できる",
			args: args{
				value: 1,
			},
			want:    value.Size(1),
			wantErr: false,
		},
		{
			name: "100で作成できる",
			args: args{
				value: 100,
			},
			want:    value.Size(100),
			wantErr: false,
		},
		{
			name: "0で作成できない",
			args: args{
				value: 0,
			},
			want:    value.Size(-1),
			wantErr: true,
		},
		{
			name: "101で作成できない",
			args: args{
				value: 101,
			},
			want:    value.Size(-1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := value.NewSize(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSize() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("NewSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
