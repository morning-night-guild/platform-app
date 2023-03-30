package auth_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"reflect"
	"strings"
	"testing"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestDecodePublicKey(t *testing.T) {
	t.Parallel()

	type args struct {
		key *Key
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "デコードできる",
			args: args{
				key: NewKey(t),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := auth.DecodePublicKey(tt.args.key.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodePublicKey() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.args.key.Key.PublicKey) {
				t.Errorf("DecodePublicKey() = %v, want %v", got, tt.args.key.Key.PublicKey)
			}
		})
	}
}

type Key struct {
	T   *testing.T
	Key *rsa.PrivateKey
}

func NewKey(t *testing.T) *Key {
	t.Helper()

	prv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	return &Key{
		T:   t,
		Key: prv,
	}
}

func (k *Key) String() string {
	k.T.Helper()

	pubkey, err := x509.MarshalPKIXPublicKey(&k.Key.PublicKey)
	if err != nil {
		k.T.Fatal(err)
	}

	b := new(bytes.Buffer)

	if err := pem.Encode(b, &pem.Block{
		Bytes: pubkey,
	}); err != nil {
		k.T.Fatal(err)
	}

	remove := func(arr []string, i int) []string {
		return append(arr[:i], arr[i+1:]...)
	}

	pems := strings.Split(b.String(), "\n")

	pems = remove(pems, len(pems)-1)

	pems = remove(pems, len(pems)-1)

	pems = remove(pems, 0)

	return strings.Join(pems, "")
}
