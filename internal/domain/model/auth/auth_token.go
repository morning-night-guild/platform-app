package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

const AuthTokenKey = "auth-token"

type AuthToken string //nolint:revive

func NewAuthToken(
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
) AuthToken {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": userID.String(),
		"iat": now.Unix(),
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

func (at AuthToken) GetUserID(secret Secret) user.ID {
	parsedToken, _ := jwt.Parse(at.String(), func(token *jwt.Token) (interface{}, error) {
		return []byte(secret.String()), nil
	})

	claims, _ := parsedToken.Claims.(jwt.MapClaims)

	id, _ := claims["sub"].(string)

	uid, _ := user.NewID(id)

	return uid
}
