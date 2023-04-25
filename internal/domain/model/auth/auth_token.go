package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

const AuthTokenKey = "auth-token"

type AuthToken string //nolint:revive

func ParseAuthToken(
	token string,
	secret Secret,
) (AuthToken, error) {
	at := AuthToken(token)

	if err := at.validate(secret); err != nil {
		return AuthToken(""), err
	}

	return at, nil
}

func GenerateAuthToken(
	userID user.ID,
	secret Secret,
	expiresIn ExpiresIn,
) AuthToken {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": userID.String(),
		"iat": now.Unix(),
		"exp": now.Add(expiresIn.Duration()).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	strToken, _ := token.SignedString([]byte(secret))

	return AuthToken(strToken)
}

func (at AuthToken) validate(secret Secret) error {
	parsedToken, err := jwt.Parse(at.String(), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])

			return nil, errors.NewValidationError(msg)
		}

		return []byte(secret.String()), nil
	})
	if err != nil {
		return err
	}

	if !parsedToken.Valid {
		return errors.NewValidationError("invalid token")
	}

	return nil
}

func (at AuthToken) String() string {
	return string(at)
}

func (at AuthToken) UserID() user.ID {
	decoded := strings.Split(at.String(), ".")

	dec, err := base64.RawStdEncoding.Strict().DecodeString(decoded[1])
	if err != nil {
		return user.GenerateZeroID()
	}

	type payload struct {
		Sub string `json:"sub"`
	}

	var p payload

	if err := json.Unmarshal(dec, &p); err != nil {
		return user.GenerateZeroID()
	}

	uid, err := user.NewID(p.Sub)
	if err != nil {
		return user.GenerateZeroID()
	}

	return uid
}
