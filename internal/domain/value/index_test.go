package value_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

func TestNewIndex(t *testing.T) {
	t.Parallel()

	type args struct {
		value int
	}

	tests := []struct {
		name    string
		args    args
		want    value.Index
		wantErr bool
	}{
		{
			name: "0で作成できる",
			args: args{
				value: 0,
			},
			want:    value.Index(0),
			wantErr: false,
		},
		{
			name: "1で作成できる",
			args: args{
				value: 1,
			},
			want:    value.Index(1),
			wantErr: false,
		},
		{
			name: "-1で作成できない",
			args: args{
				value: -1,
			},
			want:    value.Index(-1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := value.NewIndex(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIndex() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("NewIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
