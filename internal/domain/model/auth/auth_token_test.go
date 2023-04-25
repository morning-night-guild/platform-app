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

func TestParseAuthTokenToken(t *testing.T) {
	t.Parallel()

	type args struct {
		userID        user.ID
		encryptSecret auth.Secret
		expiresIn     auth.ExpiresIn
		decryptSecret auth.Secret
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "認証トークンが生成できる",
			args: args{
				userID:        user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				encryptSecret: auth.Secret("secret"),
				expiresIn:     auth.DefaultExpiresIn,
				decryptSecret: auth.Secret("secret"),
			},
			wantErr: false,
		},
		{
			name: "認証トークンの有効期限が切れている",
			args: args{
				userID:        user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				encryptSecret: auth.Secret("secret"),
				expiresIn:     auth.ExpiresIn(0),
				decryptSecret: auth.Secret("secret"),
			},
			wantErr: true,
		},
		{
			name: "暗号化と復号化のシークレットが異なる",
			args: args{
				userID:        user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				encryptSecret: auth.Secret("sencrypt"),
				expiresIn:     auth.DefaultExpiresIn,
				decryptSecret: auth.Secret("decrypt"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := auth.GenerateAuthToken(tt.args.userID, tt.args.encryptSecret, tt.args.expiresIn)
			if _, err := auth.ParseAuthToken(got.String(), tt.args.decryptSecret); (err != nil) != tt.wantErr {
				t.Errorf("ParseAuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
