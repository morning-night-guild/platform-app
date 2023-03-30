package helper

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

type User struct {
	T        *testing.T
	EMail    string
	Password string
	Cookies  []*http.Cookie
	Client   *OpenAPIClient
	Key      *rsa.PrivateKey
}

func NewUser(
	t *testing.T,
	url string,
) User {
	t.Helper()

	client := NewOpenAPIClient(t, url)

	id := uuid.New().String()

	email := fmt.Sprintf("%s@example.com", id)

	password := id

	if _, err := client.Client.V1AuthSignUp(context.Background(), openapi.V1AuthSignUpJSONRequestBody{
		Email:    types.Email(email),
		Password: password,
	}); err != nil {
		t.Fatalf("failed to auth sign up: %s", err)
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Client.V1AuthSignIn(context.Background(), openapi.V1AuthSignInJSONRequestBody{
		Email:     types.Email(email),
		Password:  password,
		PublicKey: Public(t, priv),
	})
	if err != nil {
		t.Fatalf("failed to auth sign in: %s", err)
	}

	client.Client.Client = &http.Client{
		Transport: NewCookiesTransport(t, res.Cookies()),
	}

	return User{
		T:        t,
		EMail:    email,
		Password: password,
		Cookies:  res.Cookies(),
		Client:   client,
		Key:      priv,
	}
}

func Public(t *testing.T, private *rsa.PrivateKey) string {
	t.Helper()

	b := new(bytes.Buffer)

	bt, err := x509.MarshalPKIXPublicKey(&private.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	pem.Encode(b, &pem.Block{
		Bytes: bt,
	})

	remove := func(arr []string, i int) []string {
		return append(arr[:i], arr[i+1:]...)
	}

	pems := strings.Split(b.String(), "\n")

	pems = remove(pems, len(pems)-1)

	pems = remove(pems, len(pems)-1)

	pems = remove(pems, 0)

	return strings.Join(pems, "")
}

func (user *User) Sign(code string) string {
	user.T.Helper()

	h := crypto.Hash.New(crypto.SHA256)

	h.Write([]byte(code))

	digest := h.Sum(nil)

	signed, err := rsa.SignPSS(rand.Reader, user.Key, crypto.SHA256, digest, nil)
	if err != nil {
		user.T.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(signed)
}

func ExtractUserID(t *testing.T, cookies []*http.Cookie) string {
	t.Helper()

	var cookie *http.Cookie

	for _, c := range cookies {
		if c.Name == auth.AuthTokenKey {
			cookie = c
			break
		}
	}

	log.Printf("cookie: %+v", cookie)

	payload := strings.Split(cookie.Value, ".")[1]

	decode, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		t.Fatal(err)
	}

	type token struct {
		Sub string `json:"sub"`
	}

	var p token

	if err := json.Unmarshal(decode, &p); err != nil {
		t.Fatal(err)
	}

	return p.Sub
}
