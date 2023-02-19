package article_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
)

func TestNewTag(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    article.Tag
		wantErr bool
	}{
		{
			name: "空でないタグが作成できる",
			args: args{
				value: "タイトル",
			},
			want:    article.Tag("タイトル"),
			wantErr: false,
		},
		{
			name: "空のタグは作成に失敗する",
			args: args{
				value: "",
			},
			want:    article.Tag(""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := article.NewTag(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTag() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("NewTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
