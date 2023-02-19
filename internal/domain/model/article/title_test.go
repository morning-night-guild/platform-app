package article_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
)

func TestNewTitle(t *testing.T) {
	t.Parallel()

	type args struct {
		t string
	}

	tests := []struct {
		name    string
		args    args
		want    article.Title
		wantErr bool
	}{
		{
			name: "空のタイトルが作成できる",
			args: args{
				t: "",
			},
			want:    article.Title(""),
			wantErr: false,
		},
		{
			name: "空でないタイトルが作成できる",
			args: args{
				t: "タイトル",
			},
			want:    article.Title("タイトル"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := article.NewTitle(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTitle() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("NewTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
