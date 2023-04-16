package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/rpc"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ usecase.APIAuth = (*APIAuth)(nil)

type APIAuth struct {
	authRPC      rpc.Auth
	userRPC      rpc.User
	authCache    cache.Cache[model.Auth]
	codeCache    cache.Cache[model.Code]
	sessionCache cache.Cache[model.Session]
}

func NewAPIAuth(
	authRPC rpc.Auth,
	userRPC rpc.User,
	authCache cache.Cache[model.Auth],
	codeCache cache.Cache[model.Code],
	sessionCache cache.Cache[model.Session],
) *APIAuth {
	return &APIAuth{
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

	sCmd, err := aa.sessionCache.CreateTxSetCmd(ctx, session.SessionID.String(), session, model.DefaultSessionExpiresIn)
	if err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	at := model.IssueAuth(user.UserID)

	aCmd, err := aa.authCache.CreateTxSetCmd(ctx, at.UserID.String(), at, model.DefaultAuthExpiresIn)
	if err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	if err := aa.sessionCache.Tx(ctx, []cache.TxSetCmd{sCmd, aCmd}, []cache.TxDelCmd{}); err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	return usecase.APIAuthSignInOutput{
		AuthToken:    at.ToToken(session.SessionID.ToSecret()),
		SessionToken: session.ToToken(input.Secret),
	}, nil
}

func (aa *APIAuth) SignOut(
	ctx context.Context,
	input usecase.APIAuthSignOutInput,
) (usecase.APIAuthSignOutOutput, error) {
	if err := aa.sessionCache.Del(ctx, input.SessionID.String()); err != nil {
		log.GetLogCtx(ctx).Warn("failed to delete session cache", log.ErrorField(err))
	}

	if err := aa.authCache.Del(ctx, input.UserID.String()); err != nil {
		log.GetLogCtx(ctx).Warn("failed to delete auth cache", log.ErrorField(err))
	}

	return usecase.APIAuthSignOutOutput{}, nil
}

func (aa *APIAuth) Verify(
	ctx context.Context,
	input usecase.APIAuthVerifyInput,
) (usecase.APIAuthVerifyOutput, error) {
	auth, err := aa.authCache.Get(ctx, input.UserID.String())
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
	code, err := aa.codeCache.Get(ctx, input.SessionID.String())
	if err != nil {
		return usecase.APIAuthRefreshOutput{}, errors.NewNotFoundError("code is not found", err)
	}

	if code.CodeID != input.CodeID {
		return usecase.APIAuthRefreshOutput{}, errors.NewValidationError("CodeID is invalid")
	}

	if code.IsExpired() {
		return usecase.APIAuthRefreshOutput{}, errors.NewValidationError("Code is expired")
	}

	session, err := aa.sessionCache.Get(ctx, input.SessionID.String())
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

	cCmd, err := aa.codeCache.CreateTxDelCmd(ctx, input.SessionID.String())
	if err != nil {
		return usecase.APIAuthRefreshOutput{}, err
	}

	at := model.IssueAuth(session.UserID)

	aCmd, err := aa.authCache.CreateTxSetCmd(ctx, at.UserID.String(), at, model.DefaultAuthExpiresIn)
	if err != nil {
		return usecase.APIAuthRefreshOutput{}, err
	}

	if err := aa.authCache.Tx(ctx, []cache.TxSetCmd{aCmd}, []cache.TxDelCmd{cCmd}); err != nil {
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
	code := model.GenerateCode(input.SessionID)

	if err := aa.codeCache.Set(ctx, input.SessionID.String(), code, model.DefaultCodeExpiresIn); err != nil {
		return usecase.APIAuthGenerateCodeOutput{}, err
	}

	return usecase.APIAuthGenerateCodeOutput{
		Code: code,
	}, nil
}
