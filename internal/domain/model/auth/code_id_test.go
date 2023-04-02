package auth_test

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
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

func TestCodeIDVerify(t *testing.T) {
	t.Parallel()

	type args struct {
		signature  auth.Signature
		privateKey *rsa.PrivateKey
	}

	key := generateKey(t)

	tests := []struct {
		name    string
		cd      auth.CodeID
		args    args
		wantErr bool
	}{
		{
			name: "署名が検証できる",
			cd:   auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
			args: args{
				signature:  sign(t, key, "01234567-0123-0123-0123-0123456789ab"),
				privateKey: key,
			},
			wantErr: false,
		},
		{
			name: "署名が検証できない",
			cd:   auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
			args: args{
				signature:  auth.Signature(""),
				privateKey: key,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.cd.Verify(tt.args.signature, &tt.args.privateKey.PublicKey); (err != nil) != tt.wantErr {
				t.Errorf("CodeID.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		t.Fatal(err)
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
