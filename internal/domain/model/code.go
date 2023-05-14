package model

import (
	"time"

	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
)

const DefaultCodeExpiresIn = 10 * time.Minute // 10 minutes

type Code struct {
	CodeID    auth.CodeID    `json:"codeId"`
	SessionID auth.SessionID `json:"sessionId"`
	IssuedAt  time.Time      `json:"issuedAt"`
	ExpiresAt time.Time      `json:"expiresAt"`
}

func NewCode(
	codeID auth.CodeID,
	sessionID auth.SessionID,
	issuedAt time.Time,
	expiresAt time.Time,
) (Code, error) {
	cd := Code{
		CodeID:    codeID,
		SessionID: sessionID,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	if err := cd.validate(); err != nil {
		return Code{}, err
	}

	return cd, nil
}

func GenerateCode(
	sessionID auth.SessionID,
) Code {
	now := time.Now()

	return Code{
		CodeID:    auth.GenerateCodeID(),
		SessionID: sessionID,
		IssuedAt:  now,
		ExpiresAt: now.Add(DefaultCodeExpiresIn),
	}
}

func (cd Code) validate() error {
	if cd.ExpiresAt.Before(cd.IssuedAt) {
		return errors.NewValidationError("ExpiresAt must be greater than IssuedAt")
	}

	return nil
}

func (cd Code) IsExpired() bool {
	return cd.ExpiresAt.Before(time.Now())
}
