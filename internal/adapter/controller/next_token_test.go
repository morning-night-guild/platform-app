package controller_test

import (
	"reflect"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/adapter/controller"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
)

func TestNextTokenToIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tr   controller.NextToken
		want repository.Index
	}{
		{
			name: "トークンからインデックスを生成できる",
			tr:   controller.CreateNextTokenFromIndex(repository.Index(0)),
			want: repository.Index(0),
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
		size repository.Size
	}

	tests := []struct {
		name string
		tr   controller.NextToken
		args args
		want controller.NextToken
	}{
		{
			name: "ネクストトークンが作成できる",
			tr:   controller.CreateNextTokenFromIndex(repository.Index(0)),
			args: args{
				size: repository.Size(1),
			},
			want: controller.CreateNextTokenFromIndex(repository.Index(1)),
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
