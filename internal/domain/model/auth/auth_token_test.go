package auth_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

func TestAuthTokenUserID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		at   auth.AuthToken
		want user.ID
	}{
		{
			name: "認証トークンからUserIDを取得できる",
			at:   auth.GenerateAuthToken(user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")), auth.Secret("secret"), auth.DefaultExpiresIn),
			want: user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.at.UserID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthToken.UserID() = %v, want %v", got, tt.want)
			}
		})
	}
}
