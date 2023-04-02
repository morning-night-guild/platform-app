package article_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
)

func TestNewDescription(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    article.Description
		wantErr bool
	}{
		{
			name: "記事の説明が作成できる",
			args: args{
				value: "説明",
			},
			want:    article.Description("説明"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := article.NewDescription(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}
