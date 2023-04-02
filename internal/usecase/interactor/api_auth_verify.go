package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ port.APIAuthVerify = (*APIAuthVerify)(nil)

type APIAuthVerify struct {
	secret    auth.Secret
	authCache cache.Cache[model.Auth]
}

func NewAPIAuthVerify(
	secret auth.Secret,
	authCache cache.Cache[model.Auth],
) *APIAuthVerify {
	return &APIAuthVerify{
		secret:    secret,
		authCache: authCache,
	}
}

func (aav *APIAuthVerify) Execute(
	ctx context.Context,
	input port.APIAuthVerifyInput,
) (port.APIAuthVerifyOutput, error) {
	sid := input.SessionToken.GetID(aav.secret)

	uid := input.AuthToken.GetUserID(sid.ToSecret())

	auth, err := aav.authCache.Get(ctx, uid.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get auth cache", log.ErrorField(err))

		return port.APIAuthVerifyOutput{}, err
	}

	if auth.IsExpired() {
		return port.APIAuthVerifyOutput{}, errors.NewValidationError("auth token is expired")
	}

	return port.APIAuthVerifyOutput{}, nil
}
