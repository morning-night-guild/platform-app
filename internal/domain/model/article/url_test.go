package article_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
)

func TestNewURL(t *testing.T) {
	t.Parallel()

	type args struct {
		u string
	}

	tests := []struct {
		name    string
		args    args
		want    article.URL
		wantErr bool
	}{
		{
			name: "空でなく`https://`から始まるURLが正常に作成できる",
			args: args{
				u: "https://example.com/image",
			},
			want:    article.URL("https://example.com/image"),
			wantErr: false,
		},
		{
			name: "空でなく`https://`から始まらないURLは作成に失敗する",
			args: args{
				u: "http://example.com/image",
			},
			want:    article.URL(""),
			wantErr: true,
		},
		{
			name: "空の場合はURLの作成に失敗する",
			args: args{
				u: "",
			},
			want:    article.URL(""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := article.NewURL(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewURL() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("NewURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
