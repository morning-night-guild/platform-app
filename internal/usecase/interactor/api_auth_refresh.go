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
	sessionCache cache.Cache[model.Session]
}

func NewAPIAuthRefresh(
	secret auth.Secret,
	codeCache cache.Cache[model.Code],
	sessionCache cache.Cache[model.Session],
) APIAuthRefresh {
	return APIAuthRefresh{
		secret:       secret,
		codeCache:    codeCache,
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
		return port.APIAuthRefreshOutput{}, err
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
		return port.APIAuthRefreshOutput{}, err
	}

	if err := aar.codeCache.Del(ctx, sid.String()); err != nil {
		log.GetLogCtx(ctx).Warn(err.Error())
	}

	authToken := auth.GenerateAuthToken(session.UserID, session.SessionID.ToSecret())

	return port.APIAuthRefreshOutput{
		AuthToken: authToken,
	}, nil
}
