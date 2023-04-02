package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

const SessionTokenKey = "session-token"

type SessionToken string

func NewSessionToken(
	token string,
	secret Secret,
) (SessionToken, error) {
	st := SessionToken(token)

	if err := st.validate(secret); err != nil {
		return SessionToken(""), err
	}

	return st, nil
}

func GenerateSessionToken(
	sessionID SessionID,
	secret Secret,
) SessionToken {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": sessionID.String(),
		"iat": now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	strToken, _ := token.SignedString([]byte(secret))

	return SessionToken(strToken)
}

func (st SessionToken) validate(secret Secret) error {
	parsedToken, err := jwt.Parse(st.String(), func(token *jwt.Token) (interface{}, error) {
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

func (st SessionToken) String() string {
	return string(st)
}

func (st SessionToken) GetID(secret Secret) SessionID {
	parsedToken, _ := jwt.Parse(st.String(), func(token *jwt.Token) (interface{}, error) {
		return []byte(secret.String()), nil
	})

	claims, _ := parsedToken.Claims.(jwt.MapClaims)

	id, _ := claims["sub"].(string)

	sid, _ := NewSessionID(id)

	return sid
}

func (st SessionToken) ToSecret(secret Secret) Secret {
	return st.GetID(secret).ToSecret()
}
