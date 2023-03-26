package model

import (
	"time"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

const DefaultAuthExpiresIn = time.Hour // 1 hour

type Auth struct {
	AuthID    user.UserID `json:"authId"`
	UserID    user.UserID `json:"userId"`
	IssuedAt  time.Time   `json:"issuedAt"`
	ExpiresAt time.Time   `json:"expiresAt"`
}

func NewAuth(
	authID user.UserID,
	userID user.UserID,
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
	userID user.UserID,
) Auth {
	now := time.Now()

	return Auth{
		AuthID:    userID,
		UserID:    userID,
		IssuedAt:  now,
		ExpiresAt: now.Add(DefaultAuthExpiresIn),
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
