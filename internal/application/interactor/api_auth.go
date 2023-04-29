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

func (itr *APIAuth) SignUp(
	ctx context.Context,
	input usecase.APIAuthSignUpInput,
) (usecase.APIAuthSignUpOutput, error) {
	user, err := itr.userRPC.Create(ctx)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to create user", log.ErrorField(err))

		return usecase.APIAuthSignUpOutput{}, err
	}

	if err := itr.authRPC.SignUp(ctx, user.UserID, input.EMail, input.Password); err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign up", log.ErrorField(err))

		return usecase.APIAuthSignUpOutput{}, err
	}

	return usecase.APIAuthSignUpOutput{}, nil
}

func (itr *APIAuth) SignIn(
	ctx context.Context,
	input usecase.APIAuthSignInInput,
) (usecase.APIAuthSignInOutput, error) {
	user, err := itr.authRPC.SignIn(ctx, input.EMail, input.Password)
	if err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	session := model.IssueSession(user.UserID, input.PublicKey)

	sCmd, err := itr.sessionCache.CreateTxSetCmd(ctx, session.SessionID.String(), session, model.DefaultSessionExpiresIn)
	if err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	at := model.IssueAuth(user.UserID, input.ExpiresIn)

	aCmd, err := itr.authCache.CreateTxSetCmd(ctx, at.UserID.String(), at, at.ExpiresIn().Duration())
	if err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	if err := itr.sessionCache.Tx(ctx, []cache.TxSetCmd{sCmd, aCmd}, []cache.TxDelCmd{}); err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	return usecase.APIAuthSignInOutput{
		Auth:         at,
		AuthToken:    at.ToToken(session.SessionID.ToSecret()), // secret を model.Auth{} にトークンに変換するときに secret が不要になる
		SessionToken: session.ToToken(input.Secret),
	}, nil
}

func (itr *APIAuth) SignOut(
	ctx context.Context,
	input usecase.APIAuthSignOutInput,
) (usecase.APIAuthSignOutOutput, error) {
	if err := itr.sessionCache.Del(ctx, input.SessionID.String()); err != nil {
		log.GetLogCtx(ctx).Warn("failed to delete session cache", log.ErrorField(err))
	}

	if err := itr.authCache.Del(ctx, input.UserID.String()); err != nil {
		log.GetLogCtx(ctx).Warn("failed to delete auth cache", log.ErrorField(err))
	}

	return usecase.APIAuthSignOutOutput{}, nil
}

func (itr *APIAuth) Verify(
	ctx context.Context,
	input usecase.APIAuthVerifyInput,
) (usecase.APIAuthVerifyOutput, error) {
	auth, err := itr.authCache.Get(ctx, input.UserID.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get auth cache", log.ErrorField(err))

		return usecase.APIAuthVerifyOutput{}, errors.NewUnauthorizedError("not found auth", err)
	}

	if auth.IsExpired() {
		return usecase.APIAuthVerifyOutput{}, errors.NewUnauthorizedError("auth is expired")
	}

	return usecase.APIAuthVerifyOutput{}, nil
}

func (itr *APIAuth) Refresh(
	ctx context.Context,
	input usecase.APIAuthRefreshInput,
) (usecase.APIAuthRefreshOutput, error) {
	code, err := itr.codeCache.Get(ctx, input.SessionID.String())
	if err != nil {
		return usecase.APIAuthRefreshOutput{}, errors.NewNotFoundError("code is not found", err)
	}

	if code.CodeID != input.CodeID {
		return usecase.APIAuthRefreshOutput{}, errors.NewValidationError("CodeID is invalid")
	}

	if code.IsExpired() {
		return usecase.APIAuthRefreshOutput{}, errors.NewValidationError("Code is expired")
	}

	session, err := itr.sessionCache.Get(ctx, input.SessionID.String())
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

	cCmd, err := itr.codeCache.CreateTxDelCmd(ctx, input.SessionID.String())
	if err != nil {
		return usecase.APIAuthRefreshOutput{}, err
	}

	at := model.IssueAuth(session.UserID, input.ExpiresIn)

	aCmd, err := itr.authCache.CreateTxSetCmd(ctx, at.UserID.String(), at, model.DefaultAuthExpiresIn)
	if err != nil {
		return usecase.APIAuthRefreshOutput{}, err
	}

	if err := itr.authCache.Tx(ctx, []cache.TxSetCmd{aCmd}, []cache.TxDelCmd{cCmd}); err != nil {
		return usecase.APIAuthRefreshOutput{}, err
	}

	return usecase.APIAuthRefreshOutput{
		AuthToken: at.ToToken(session.SessionID.ToSecret()),
	}, nil
}

func (itr *APIAuth) GenerateCode(
	ctx context.Context,
	input usecase.APIAuthGenerateCodeInput,
) (usecase.APIAuthGenerateCodeOutput, error) {
	code := model.GenerateCode(input.SessionID)

	if err := itr.codeCache.Set(ctx, input.SessionID.String(), code, model.DefaultCodeExpiresIn); err != nil {
		return usecase.APIAuthGenerateCodeOutput{}, err
	}

	return usecase.APIAuthGenerateCodeOutput{
		Code: code,
	}, nil
}
