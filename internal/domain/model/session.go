package model

import (
	"crypto/rsa"
	"time"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

const DefaultSessionExpiresIn = time.Hour * 24 * 30

type Session struct {
	SessionID auth.SessionID `json:"sessionId"`
	UserID    user.UserID    `json:"userId"`
	PublicKey rsa.PublicKey  `json:"publicKey"`
	IssuedAt  time.Time      `json:"issuedAt"`
	ExpiresAt time.Time      `json:"expiresAt"`
}

func NewSession(
	sessionID auth.SessionID,
	userID user.UserID,
	publicKey rsa.PublicKey,
	issuedAt time.Time,
	expiresAt time.Time,
) (Session, error) {
	sss := Session{
		SessionID: sessionID,
		UserID:    userID,
		PublicKey: publicKey,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	if err := sss.validate(); err != nil {
		return Session{}, err
	}

	return sss, nil
}

func IssueSession(
	userID user.UserID,
	publicKey rsa.PublicKey,
) Session {
	now := time.Now()

	return Session{
		SessionID: auth.GenerateSessionID(),
		UserID:    userID,
		PublicKey: publicKey,
		IssuedAt:  now,
		ExpiresAt: now.Add(DefaultSessionExpiresIn),
	}
}

func (sss Session) validate() error {
	if sss.ExpiresAt.Before(sss.IssuedAt) {
		return errors.NewValidationError("ExpiresAt must be after IssuedAt")
	}

	return nil
}

func (sss Session) IsExpired() bool {
	return sss.ExpiresAt.Before(time.Now())
}

func (sss Session) ToToken(
	secret auth.Secret,
) auth.SessionToken {
	return auth.GenerateSessionToken(sss.SessionID, secret)
}
