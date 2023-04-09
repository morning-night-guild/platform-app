package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ usecase.APIAuth = (*APIAuth)(nil)

type APIAuth struct {
	secret       auth.Secret
	authRPC      rpc.Auth
	userRPC      rpc.User
	authCache    cache.Cache[model.Auth]
	codeCache    cache.Cache[model.Code]
	sessionCache cache.Cache[model.Session]
}

func NewAPIAuth(
	secret auth.Secret,
	authRPC rpc.Auth,
	userRPC rpc.User,
	authCache cache.Cache[model.Auth],
	codeCache cache.Cache[model.Code],
	sessionCache cache.Cache[model.Session],
) *APIAuth {
	return &APIAuth{
		secret:       secret,
		authRPC:      authRPC,
		userRPC:      userRPC,
		authCache:    authCache,
		codeCache:    codeCache,
		sessionCache: sessionCache,
	}
}

func (aa *APIAuth) SignUp(
	ctx context.Context,
	input usecase.APIAuthSignUpInput,
) (usecase.APIAuthSignUpOutput, error) {
	user, err := aa.userRPC.Create(ctx)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to create user", log.ErrorField(err))

		return usecase.APIAuthSignUpOutput{}, err
	}

	if err := aa.authRPC.SignUp(ctx, user.UserID, input.EMail, input.Password); err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign up", log.ErrorField(err))

		return usecase.APIAuthSignUpOutput{}, err
	}

	return usecase.APIAuthSignUpOutput{}, nil
}

func (aa *APIAuth) SignIn(
	ctx context.Context,
	input usecase.APIAuthSignInInput,
) (usecase.APIAuthSignInOutput, error) {
	user, err := aa.authRPC.SignIn(ctx, input.EMail, input.Password)
	if err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	session := model.IssueSession(user.UserID, input.PublicKey)

	if err := aa.sessionCache.Set(ctx, session.SessionID.String(), session, model.DefaultSessionExpiresIn); err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	at := model.IssueAuth(user.UserID)

	if err := aa.authCache.Set(ctx, at.UserID.String(), at, model.DefaultAuthExpiresIn); err != nil {
		if err := aa.sessionCache.Del(ctx, session.SessionID.String()); err != nil {
			log.GetLogCtx(ctx).Warn("failed to delete session", log.ErrorField(err))
		}

		return usecase.APIAuthSignInOutput{}, err
	}

	return usecase.APIAuthSignInOutput{
		AuthToken:    at.ToToken(session.SessionID.ToSecret()),
		SessionToken: session.ToToken(aa.secret),
	}, nil
}

func (aa *APIAuth) SignOut(
	ctx context.Context,
	input usecase.APIAuthSignOutInput,
) (usecase.APIAuthSignOutOutput, error) {
	sid := input.SessionToken.ID(aa.secret)

	uid := input.AuthToken.UserID(sid.ToSecret())

	if err := aa.sessionCache.Del(ctx, sid.String()); err != nil {
		log.GetLogCtx(ctx).Warn("failed to delete session cache", log.ErrorField(err))
	}

	if err := aa.authCache.Del(ctx, uid.String()); err != nil {
		log.GetLogCtx(ctx).Warn("failed to delete auth cache", log.ErrorField(err))
	}

	return usecase.APIAuthSignOutOutput{}, nil
}

func (aa *APIAuth) Verify(
	ctx context.Context,
	input usecase.APIAuthVerifyInput,
) (usecase.APIAuthVerifyOutput, error) {
	sid := input.SessionToken.ID(aa.secret)

	uid := input.AuthToken.UserID(sid.ToSecret())

	auth, err := aa.authCache.Get(ctx, uid.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get auth cache", log.ErrorField(err))

		return usecase.APIAuthVerifyOutput{}, err
	}

	if auth.IsExpired() {
		return usecase.APIAuthVerifyOutput{}, errors.NewUnauthorizedError("auth token is expired")
	}

	return usecase.APIAuthVerifyOutput{}, nil
}

func (aa *APIAuth) Refresh(
	ctx context.Context,
	input usecase.APIAuthRefreshInput,
) (usecase.APIAuthRefreshOutput, error) {
	sid := input.SessionToken.ID(aa.secret)

	code, err := aa.codeCache.Get(ctx, sid.String())
	if err != nil {
		return usecase.APIAuthRefreshOutput{}, errors.NewNotFoundError("code is not found", err)
	}

	if code.CodeID != input.CodeID {
		return usecase.APIAuthRefreshOutput{}, errors.NewValidationError("CodeID is invalid")
	}

	if code.IsExpired() {
		return usecase.APIAuthRefreshOutput{}, errors.NewValidationError("Code is expired")
	}

	session, err := aa.sessionCache.Get(ctx, sid.String())
	if err != nil {
		return usecase.APIAuthRefreshOutput{}, err
	}

	if session.IsExpired() {
		return usecase.APIAuthRefreshOutput{}, errors.NewValidationError("Session is expired")
	}

	if err := code.CodeID.Verify(input.Signature, &session.PublicKey); err != nil {
		log.GetLogCtx(ctx).Warn("signature invalid", log.ErrorField(err))

		return usecase.APIAuthRefreshOutput{}, errors.NewUnauthorizedError("signature invalid")
	}

	if err := aa.codeCache.Del(ctx, sid.String()); err != nil {
		log.GetLogCtx(ctx).Warn(err.Error())
	}

	at := model.IssueAuth(session.UserID)

	if err := aa.authCache.Set(ctx, at.UserID.String(), at, model.DefaultAuthExpiresIn); err != nil {
		if err := aa.sessionCache.Del(ctx, session.SessionID.String()); err != nil {
			log.GetLogCtx(ctx).Warn("failed to delete session", log.ErrorField(err))
		}

		return usecase.APIAuthRefreshOutput{}, err
	}

	return usecase.APIAuthRefreshOutput{
		AuthToken: at.ToToken(session.SessionID.ToSecret()),
	}, nil
}

func (aa *APIAuth) GenerateCode(
	ctx context.Context,
	input usecase.APIAuthGenerateCodeInput,
) (usecase.APIAuthGenerateCodeOutput, error) {
	sid := input.SessionToken.ID(aa.secret)

	code := model.GenerateCode(sid)

	if err := aa.codeCache.Set(ctx, sid.String(), code, model.DefaultCodeExpiresIn); err != nil {
		return usecase.APIAuthGenerateCodeOutput{}, err
	}

	return usecase.APIAuthGenerateCodeOutput{
		Code: code,
	}, nil
}
