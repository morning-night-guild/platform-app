package value_test

import (
	"reflect"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/value"
)

func TestNextTokenToIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tr   value.NextToken
		want value.Index
	}{
		{
			name: "トークンからインデックスを生成できる",
			tr:   value.CreateNextTokenFromIndex(value.Index(0)),
			want: value.Index(0),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.tr.ToIndex(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NextToken.ToIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextTokenCreateNextNextToken(t *testing.T) {
	t.Parallel()

	type args struct {
		size value.Size
	}

	tests := []struct {
		name string
		tr   value.NextToken
		args args
		want value.NextToken
	}{
		{
			name: "ネクストトークンが作成できる",
			tr:   value.CreateNextTokenFromIndex(value.Index(0)),
			args: args{
				size: value.Size(1),
			},
			want: value.CreateNextTokenFromIndex(value.Index(1)),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.tr.CreateNextToken(tt.args.size); got != tt.want {
				t.Errorf("NextToken.CreateNextNextToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
