package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ port.APIAuthSignOut = (*APIAuthSignOut)(nil)

type APIAuthSignOut struct {
	secret       auth.Secret
	authCache    cache.Cache[model.Auth]
	sessionCache cache.Cache[model.Session]
}

func NewAPIAuthSignOut(
	secret auth.Secret,
	authCache cache.Cache[model.Auth],
	sessionCache cache.Cache[model.Session],
) APIAuthSignOut {
	return APIAuthSignOut{
		secret:       secret,
		authCache:    authCache,
		sessionCache: sessionCache,
	}
}

func (aas *APIAuthSignOut) Execute(
	ctx context.Context,
	input port.APIAuthSignOutInput,
) (port.APIAuthSignOutOutput, error) {
	sid := input.SessionToken.GetID(aas.secret)

	uid := input.AuthToken.GetUserID(sid.ToSecret())

	if err := aas.sessionCache.Del(ctx, sid.String()); err != nil {
		log.GetLogCtx(ctx).Warn("failed to delete session cache", log.ErrorField(err))
	}

	if err := aas.authCache.Del(ctx, uid.String()); err != nil {
		log.GetLogCtx(ctx).Warn("failed to delete auth cache", log.ErrorField(err))
	}

	return port.APIAuthSignOutOutput{}, nil
}
