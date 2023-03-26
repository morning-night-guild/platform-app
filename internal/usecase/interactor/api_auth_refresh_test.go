package interactor_test

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/usecase/interactor"
	"github.com/morning-night-guild/platform-app/internal/usecase/mock"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

func TestAPIAuthRefreshExecute(t *testing.T) {
	t.Parallel()

	type fields struct {
		secret       auth.Secret
		codeCache    cache.Cache[model.Code]
		sessionCache cache.Cache[model.Session]
	}

	type args struct {
		ctx   context.Context
		input port.APIAuthRefreshInput
	}

	key := generateKey(t)

	now := time.Now()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    port.APIAuthRefreshOutput
		wantErr bool
	}{
		{
			name: "リフレッシュできる",
			fields: fields{
				secret: auth.Secret("secret"),
				codeCache: &mock.Cache[model.Code]{
					T: t,
					Value: model.Code{
						CodeID:    auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						IssuedAt:  now,
						ExpiresAt: now.Add(model.DefaultCodeExpiresIn),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
					DelAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
				sessionCache: &mock.Cache[model.Session]{
					T: t,
					Value: model.Session{
						SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						UserID:    user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						PublicKey: key.PublicKey,
						IssuedAt:  now,
						ExpiresAt: now.Add(model.DefaultSessionExpiresIn),
					},
					GetAssert: func(t *testing.T, key string) {
						t.Helper()
					},
				},
			},
			args: args{
				ctx: context.Background(),
				input: port.APIAuthRefreshInput{
					CodeID:    auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
					Signature: sign(t, key, "01234567-0123-0123-0123-0123456789ab"),
					SessionToken: auth.GenerateSessionToken(
						auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
						auth.Secret("secret"),
					),
				},
			},
			want:    port.APIAuthRefreshOutput{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			aar := interactor.NewAPIAuthRefresh(
				tt.fields.secret,
				tt.fields.codeCache,
				tt.fields.sessionCache,
			)
			_, err := aar.Execute(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIAuthRefresh.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("APIAuthRefresh.Execute() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func sign(t *testing.T, prv *rsa.PrivateKey, code string) auth.Signature {
	t.Helper()

	h := crypto.Hash.New(crypto.SHA256)

	h.Write([]byte(code))

	hashed := h.Sum(nil)

	signed, err := rsa.SignPSS(rand.Reader, prv, crypto.SHA256, hashed, nil)
	if err != nil {
		panic(err)
	}

	return auth.Signature(base64.StdEncoding.EncodeToString(signed))
}

func generateKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	prv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	return prv
}
