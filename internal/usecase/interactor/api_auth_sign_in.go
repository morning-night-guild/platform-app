package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ port.APIAuthSignIn = (*APIAuthSignIn)(nil)

type APIAuthSignIn struct {
	secret       auth.Secret
	authRPC      rpc.Auth
	authCache    cache.Cache[model.Auth]
	sessionCache cache.Cache[model.Session]
}

func NewAPIAuthSignIn(
	secret auth.Secret,
	authRPC rpc.Auth,
	authCache cache.Cache[model.Auth],
	sessionCache cache.Cache[model.Session],
) APIAuthSignIn {
	return APIAuthSignIn{
		secret:       secret,
		authRPC:      authRPC,
		authCache:    authCache,
		sessionCache: sessionCache,
	}
}

func (aas *APIAuthSignIn) Execute(
	ctx context.Context,
	input port.APIAuthSignInInput,
) (port.APIAuthSignInOutput, error) {
	user, err := aas.authRPC.SignIn(ctx, input.EMail, input.Password)
	if err != nil {
		return port.APIAuthSignInOutput{}, err
	}

	session := model.IssueSession(user.UserID, input.PublicKey)

	if err := aas.sessionCache.Set(ctx, session.SessionID.String(), session, model.DefaultSessionExpiresIn); err != nil {
		return port.APIAuthSignInOutput{}, err
	}

	at := model.IssueAuth(user.UserID)

	if err := aas.authCache.Set(ctx, at.UserID.String(), at, model.DefaultAuthExpiresIn); err != nil {
		if err := aas.sessionCache.Del(ctx, session.SessionID.String()); err != nil {
			log.GetLogCtx(ctx).Warn("failed to delete session", log.ErrorField(err))
		}

		return port.APIAuthSignInOutput{}, err
	}

	return port.APIAuthSignInOutput{
		AuthToken:    at.ToToken(session.SessionID.ToSecret()),
		SessionToken: session.ToToken(aas.secret),
	}, nil
}
