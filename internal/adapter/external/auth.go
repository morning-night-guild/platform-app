package external

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ rpc.Auth = (*Auth)(nil)

type AuthFactory interface {
	Auth(secret string, endpoint string, apiKey string) (*Auth, error)
}

type Auth struct {
	endpoint     string
	apiKey       string
	httpClient   *http.Client
	firebaseAuth *firebase.Client
}

func NewAuth(
	endpoint string,
	apiKey string,
	httpClient *http.Client,
	firebaseAuth *firebase.Client,
) *Auth {
	return &Auth{
		endpoint:     endpoint,
		apiKey:       apiKey,
		httpClient:   httpClient,
		firebaseAuth: firebaseAuth,
	}
}

func (at *Auth) SignUp(ctx context.Context, uid user.ID, email auth.EMail, password auth.Password) error {
	params := (&firebase.UserToCreate{}).
		UID(uid.String()).
		Email(email.String()).
		EmailVerified(false).
		Password(password.String()).
		Disabled(false)

	if _, err := at.firebaseAuth.CreateUser(ctx, params); err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign up user", log.ErrorField(err))

		return err
	}

	return nil
}

type SignInRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type SignInResponse struct {
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

//nolint:funlen
func (at *Auth) SignIn(ctx context.Context, email auth.EMail, password auth.Password) (model.User, error) {
	// https://firebase.google.com/docs/reference/rest/auth#section-sign-in-email-password
	url := fmt.Sprintf("%s/v1/accounts:signInWithPassword?key=%s", at.endpoint, at.apiKey)

	req := SignInRequest{
		Email:             email.String(),
		Password:          password.String(),
		ReturnSecureToken: true,
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		log.GetLogCtx(ctx).Warn("failed to encode json", log.ErrorField(err))

		msg := fmt.Sprintf("failed to encode json. caused by %v", err)

		return model.User{}, errors.NewValidationError(msg)
	}

	res, err := at.httpClient.Post(url, "application/json", &buf) //nolint:noctx
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to post "+url, log.ErrorField(err))

		return model.User{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		message, err := io.ReadAll(res.Body)
		if err != nil {
			log.GetLogCtx(ctx).Warn("failed to read body", log.ErrorField(err))

			message = []byte(fmt.Sprintf("could not load message caused by %v", err))
		}

		msg := fmt.Sprintf("firebase error. status code is %d, message is %v", res.StatusCode, string(message))

		log.GetLogCtx(ctx).Warn(msg)

		if res.StatusCode == http.StatusUnauthorized {
			return model.User{}, errors.NewUnauthorizedError("invalid email or password")
		}

		return model.User{}, errors.NewUnknownError("unknown error")
	}

	var resp SignInResponse

	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		log.GetLogCtx(ctx).Warn("failed to decode json", log.ErrorField(err))

		return model.User{}, err
	}

	strs := strings.Split(resp.IDToken, ".")

	tmpPayload, err := base64.RawStdEncoding.DecodeString(strs[1])
	if err != nil {
		return model.User{}, fmt.Errorf("failed to decode payload: %w", err)
	}

	type Payload struct {
		UserID string `json:"user_id"` //nolint:tagliatelle
	}

	var payload Payload

	if err := json.Unmarshal(tmpPayload, &payload); err != nil {
		return model.User{}, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	uid := user.ID(uuid.MustParse(payload.UserID))

	return model.User{
		UserID: uid,
	}, nil
}
