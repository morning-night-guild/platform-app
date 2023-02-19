package article_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
)

func TestNewID(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    article.ID
		wantErr bool
	}{
		{
			name: "IDを生成できる",
			args: args{
				value: "2f8e01fb-bf67-45cc-83b0-4cfa0548a9b2",
			},
			want:    article.ID(uuid.MustParse("2f8e01fb-bf67-45cc-83b0-4cfa0548a9b2")),
			wantErr: false,
		},
		{
			name: "不正な文字列によりIDが生成できない",
			args: args{
				value: "id",
			},
			want:    article.ID{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := article.NewID(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewID() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewID() = %v, want %v", got, tt.want)
			}
		})
	}
}
