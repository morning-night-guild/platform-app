package auth_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestNewCodeID(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    auth.CodeID
		wantErr bool
	}{
		{
			name: "コードが作成できる",
			args: args{
				value: "01234567-0123-0123-0123-0123456789ab",
			},
			want:    auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
			wantErr: false,
		},
		{
			name: "コードが作成できない",
			args: args{
				value: "id",
			},
			want:    auth.CodeID{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := auth.NewCodeID(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCodeID() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCodeID() = %v, want %v", got, tt.want)
			}
		})
	}
}
