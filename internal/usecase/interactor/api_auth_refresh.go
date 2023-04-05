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

var _ port.APIAuthRefresh = (*APIAuthRefresh)(nil)

type APIAuthRefresh struct {
	secret       auth.Secret
	codeCache    cache.Cache[model.Code]
	authCache    cache.Cache[model.Auth]
	sessionCache cache.Cache[model.Session]
}

func NewAPIAuthRefresh(
	secret auth.Secret,
	codeCache cache.Cache[model.Code],
	authCache cache.Cache[model.Auth],
	sessionCache cache.Cache[model.Session],
) *APIAuthRefresh {
	return &APIAuthRefresh{
		secret:       secret,
		codeCache:    codeCache,
		authCache:    authCache,
		sessionCache: sessionCache,
	}
}

func (aar *APIAuthRefresh) Execute(
	ctx context.Context,
	input port.APIAuthRefreshInput,
) (port.APIAuthRefreshOutput, error) {
	sid := input.SessionToken.GetID(aar.secret)

	code, err := aar.codeCache.Get(ctx, sid.String())
	if err != nil {
		return port.APIAuthRefreshOutput{}, errors.NewNotFoundError("code is not found", err)
	}

	if code.CodeID != input.CodeID {
		return port.APIAuthRefreshOutput{}, errors.NewValidationError("CodeID is invalid")
	}

	if code.IsExpired() {
		return port.APIAuthRefreshOutput{}, errors.NewValidationError("Code is expired")
	}

	session, err := aar.sessionCache.Get(ctx, sid.String())
	if err != nil {
		return port.APIAuthRefreshOutput{}, err
	}

	if session.IsExpired() {
		return port.APIAuthRefreshOutput{}, errors.NewValidationError("Session is expired")
	}

	if err := code.CodeID.Verify(input.Signature, &session.PublicKey); err != nil {
		log.GetLogCtx(ctx).Warn("signature invalid", log.ErrorField(err))

		return port.APIAuthRefreshOutput{}, errors.NewUnauthorizedError("signature invalid")
	}

	if err := aar.codeCache.Del(ctx, sid.String()); err != nil {
		log.GetLogCtx(ctx).Warn(err.Error())
	}

	at := model.IssueAuth(session.UserID)

	if err := aar.authCache.Set(ctx, at.UserID.String(), at, model.DefaultAuthExpiresIn); err != nil {
		if err := aar.sessionCache.Del(ctx, session.SessionID.String()); err != nil {
			log.GetLogCtx(ctx).Warn("failed to delete session", log.ErrorField(err))
		}

		return port.APIAuthRefreshOutput{}, err
	}

	return port.APIAuthRefreshOutput{
		AuthToken: at.ToToken(session.SessionID.ToSecret()),
	}, nil
}
