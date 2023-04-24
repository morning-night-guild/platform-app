package model

import (
	"time"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

const DefaultAuthExpiresIn = time.Hour // 1 hour

type Auth struct {
	AuthID    user.ID   `json:"authId"`
	UserID    user.ID   `json:"userId"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewAuth(
	authID user.ID,
	userID user.ID,
	issuedAt time.Time,
	expiresAt time.Time,
) (Auth, error) {
	at := Auth{
		AuthID:    authID,
		UserID:    userID,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	if err := at.validate(); err != nil {
		return Auth{}, err
	}

	return at, nil
}

func IssueAuth(
	userID user.ID,
	expiresIn auth.ExpiresIn,
) Auth {
	now := time.Now()

	return Auth{
		AuthID:    userID,
		UserID:    userID,
		IssuedAt:  now,
		ExpiresAt: now.Add(expiresIn.Duration()),
	}
}

func (at Auth) validate() error {
	if at.ExpiresAt.Before(at.IssuedAt) {
		return errors.NewValidationError("ExpiresAt must be greater than IssuedAt")
	}

	return nil
}

func (at Auth) IsExpired() bool {
	return at.ExpiresAt.Before(time.Now())
}

func (at Auth) ToToken(
	secret auth.Secret,
) auth.AuthToken {
	return auth.GenerateAuthToken(at.UserID, secret)
}

func (at Auth) ExpiresIn() auth.ExpiresIn {
	expiresIn := at.ExpiresAt.Unix() - time.Now().Unix()

	if expiresIn > 0 {
		return auth.ExpiresIn(expiresIn)
	}

	return auth.ExpiresIn(0)
}
