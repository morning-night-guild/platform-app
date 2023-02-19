package repository_test

import (
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/repository"
)

func TestNewIndex(t *testing.T) {
	t.Parallel()

	type args struct {
		value int
	}

	tests := []struct {
		name    string
		args    args
		want    repository.Index
		wantErr bool
	}{
		{
			name: "0で作成できる",
			args: args{
				value: 0,
			},
			want:    repository.Index(0),
			wantErr: false,
		},
		{
			name: "1で作成できる",
			args: args{
				value: 1,
			},
			want:    repository.Index(1),
			wantErr: false,
		},
		{
			name: "-1で作成できない",
			args: args{
				value: -1,
			},
			want:    repository.Index(-1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := repository.NewIndex(tt.args.value)
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

func TestNewSize(t *testing.T) {
	t.Parallel()

	type args struct {
		value int
	}

	tests := []struct {
		name    string
		args    args
		want    repository.Size
		wantErr bool
	}{
		{
			name: "1で作成できる",
			args: args{
				value: 1,
			},
			want:    repository.Size(1),
			wantErr: false,
		},
		{
			name: "100で作成できる",
			args: args{
				value: 100,
			},
			want:    repository.Size(100),
			wantErr: false,
		},
		{
			name: "0で作成できない",
			args: args{
				value: 0,
			},
			want:    repository.Size(-1),
			wantErr: true,
		},
		{
			name: "101で作成できない",
			args: args{
				value: 101,
			},
			want:    repository.Size(-1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := repository.NewSize(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSize() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("NewSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
