package value

import (
	"reflect"
	"testing"
)

func TestNewFilter(t *testing.T) {
	t.Parallel()

	type args struct {
		name  string
		value string
	}

	tests := []struct {
		name string
		args args
		want Filter
	}{
		{
			name: "Filterを生成できる",
			args: args{
				name:  "name",
				value: "value",
			},
			want: Filter{
				Name:  "name",
				Value: "value",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewFilter(tt.args.name, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
