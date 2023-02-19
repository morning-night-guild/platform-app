package article_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
)

func TestNewThumbnail(t *testing.T) {
	t.Parallel()

	type args struct {
		t string
	}

	tests := []struct {
		name    string
		args    args
		want    article.Thumbnail
		wantErr bool
	}{
		{
			name: "空のサムネイルが作成できる",
			args: args{
				t: "",
			},
			want:    article.Thumbnail(""),
			wantErr: false,
		},
		{
			name: "空でなく`https://`から始まるサムネイルが作成できる",
			args: args{
				t: "https://example.com/image",
			},
			want:    article.Thumbnail("https://example.com/image"),
			wantErr: false,
		},
		{
			name: "空でなく`https://`から始まらないサムネイルは作成に失敗する",
			args: args{
				t: "http://example.com/image",
			},
			want:    article.Thumbnail(""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := article.NewThumbnail(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewThumbnail() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("NewThumbnail() = %v, want %v", got, tt.want)
			}
		})
	}
}
