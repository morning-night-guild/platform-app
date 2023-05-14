package interactor

import (
	"context"
	"fmt"
	"strings"

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
	userCache    cache.Cache[model.User]
	authCache    cache.Cache[model.Auth]
	codeCache    cache.Cache[model.Code]
	sessionCache cache.Cache[model.Session]
}

func NewAPIAuth(
	authRPC rpc.Auth,
	userRPC rpc.User,
	userCache cache.Cache[model.User],
	authCache cache.Cache[model.Auth],
	codeCache cache.Cache[model.Code],
	sessionCache cache.Cache[model.Session],
) *APIAuth {
	return &APIAuth{
		authRPC:      authRPC,
		userCache:    userCache,
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

	uCmd, err := itr.userCache.CreateTxSetCmd(ctx, session.SessionID.String(), user, model.DefaultSessionExpiresIn)
	if err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	sCmd, err := itr.sessionCache.CreateTxSetCmd(ctx, session.Key(), session, model.DefaultSessionExpiresIn)
	if err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	at := model.IssueAuth(user.UserID, input.ExpiresIn)

	aCmd, err := itr.authCache.CreateTxSetCmd(ctx, at.UserID.String(), at, at.ExpiresIn().Duration())
	if err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	if err := itr.sessionCache.Tx(ctx, []cache.TxSetCmd{uCmd, sCmd, aCmd}, []cache.TxDelCmd{}); err != nil {
		return usecase.APIAuthSignInOutput{}, err
	}

	return usecase.APIAuthSignInOutput{
		Auth:         at,
		AuthToken:    at.ToToken(session.SessionID.ToSecret()),
		SessionToken: session.ToToken(input.Secret),
	}, nil
}

func (itr *APIAuth) SignOut(
	ctx context.Context,
	input usecase.APIAuthSignOutInput,
) (usecase.APIAuthSignOutOutput, error) {
	key := fmt.Sprintf(model.SessionKeyFormat, input.UserID.String(), input.SessionID.String())

	sessionDelCmd, err := itr.sessionCache.CreateTxDelCmd(ctx, key)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to create session cache delete command", log.ErrorField(err))
	}

	userDelCmd, err := itr.userCache.CreateTxDelCmd(ctx, input.SessionID.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to create user cache delete command", log.ErrorField(err))
	}

	authDelCmd, err := itr.authCache.CreateTxDelCmd(ctx, input.UserID.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to create auth cache delete command", log.ErrorField(err))
	}

	delCmds := []cache.TxDelCmd{sessionDelCmd, userDelCmd, authDelCmd}

	if err := itr.sessionCache.Tx(ctx, []cache.TxSetCmd{}, delCmds); err != nil {
		log.GetLogCtx(ctx).Warn("failed to execute transaction", log.ErrorField(err))
	}

	return usecase.APIAuthSignOutOutput{}, nil
}

func (itr *APIAuth) SignOutAll(
	ctx context.Context,
	input usecase.APIAuthSignOutAllInput,
) (usecase.APIAuthSignOutAllOutput, error) {
	keys, err := itr.sessionCache.Keys(ctx, input.UserID.String(), cache.WithoutPrefix)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get session cache keys", log.ErrorField(err))
	}

	const length = 1

	delCmds := make([]cache.TxDelCmd, 0, len(keys)*2+length)

	for _, key := range keys {
		sDelCmd, err := itr.sessionCache.CreateTxDelCmd(ctx, key)
		if err != nil {
			log.GetLogCtx(ctx).Warn("failed to create session cache delete command", log.ErrorField(err))

			continue
		}

		delCmds = append(delCmds, sDelCmd)

		// 取得したkey:${user_id}:${session_id}形式 -> ここから${session_id}のみを抽出
		ukey := strings.ReplaceAll(key, fmt.Sprintf("%s:", input.UserID.String()), "")

		uDelCmd, err := itr.userCache.CreateTxDelCmd(ctx, ukey)
		if err != nil {
			log.GetLogCtx(ctx).Warn("failed to create user cache delete command", log.ErrorField(err))

			continue
		}

		delCmds = append(delCmds, uDelCmd)
	}

	authDelCmd, err := itr.authCache.CreateTxDelCmd(ctx, input.UserID.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to create auth cache delete command", log.ErrorField(err))
	} else {
		delCmds = append(delCmds, authDelCmd)
	}

	if err := itr.sessionCache.Tx(ctx, []cache.TxSetCmd{}, delCmds); err != nil {
		log.GetLogCtx(ctx).Warn("failed to transaction", log.ErrorField(err))
	}

	return usecase.APIAuthSignOutAllOutput{}, nil
}

func (itr *APIAuth) Verify(
	ctx context.Context,
	input usecase.APIAuthVerifyInput,
) (usecase.APIAuthVerifyOutput, error) {
	key := fmt.Sprintf(model.SessionKeyFormat, input.UserID.String(), input.SessionID.String())

	session, err := itr.sessionCache.Get(ctx, key)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get session cache", log.ErrorField(err))

		return usecase.APIAuthVerifyOutput{}, errors.NewUnauthorizedError("not found auth", err)
	}

	if session.IsExpired() {
		log.GetLogCtx(ctx).Warn("session is expired")

		return usecase.APIAuthVerifyOutput{}, errors.NewUnauthorizedError("session is expired")
	}

	auth, err := itr.authCache.Get(ctx, input.UserID.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get auth cache", log.ErrorField(err))

		return usecase.APIAuthVerifyOutput{}, errors.NewUnauthorizedError("not found auth", err)
	}

	if auth.IsExpired() {
		log.GetLogCtx(ctx).Warn("auth is expired")

		return usecase.APIAuthVerifyOutput{}, errors.NewUnauthorizedError("auth is expired")
	}

	return usecase.APIAuthVerifyOutput{}, nil
}

func (itr *APIAuth) Refresh( //nolint:funlen,cyclop
	ctx context.Context,
	input usecase.APIAuthRefreshInput,
) (usecase.APIAuthRefreshOutput, error) {
	code, err := itr.codeCache.Get(ctx, input.SessionID.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get code cache", log.ErrorField(err))

		return usecase.APIAuthRefreshOutput{}, errors.NewNotFoundError("code is not found", err)
	}

	if code.CodeID != input.CodeID {
		log.GetLogCtx(ctx).Warn("code id is invalid")

		return usecase.APIAuthRefreshOutput{}, errors.NewValidationError("CodeID is invalid")
	}

	if code.IsExpired() {
		log.GetLogCtx(ctx).Warn("code is expired")

		return usecase.APIAuthRefreshOutput{}, errors.NewValidationError("Code is expired")
	}

	user, err := itr.userCache.Get(ctx, input.SessionID.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get user cache", log.ErrorField(err))

		return usecase.APIAuthRefreshOutput{}, errors.NewNotFoundError("user is not found", err)
	}

	key := fmt.Sprintf(model.SessionKeyFormat, user.UserID.String(), input.SessionID.String())

	session, err := itr.sessionCache.Get(ctx, key)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get session cache", log.ErrorField(err))

		return usecase.APIAuthRefreshOutput{}, err
	}

	if session.IsExpired() {
		log.GetLogCtx(ctx).Warn("session is expired")

		return usecase.APIAuthRefreshOutput{}, errors.NewValidationError("Session is expired")
	}

	if err := code.CodeID.Verify(input.Signature, &session.PublicKey); err != nil {
		log.GetLogCtx(ctx).Warn("signature invalid", log.ErrorField(err))

		return usecase.APIAuthRefreshOutput{}, errors.NewUnauthorizedError("signature invalid")
	}

	cCmd, err := itr.codeCache.CreateTxDelCmd(ctx, input.SessionID.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to create code cache delete command", log.ErrorField(err))

		return usecase.APIAuthRefreshOutput{}, err
	}

	at := model.IssueAuth(session.UserID, input.ExpiresIn)

	aCmd, err := itr.authCache.CreateTxSetCmd(ctx, at.UserID.String(), at, model.DefaultAuthExpiresIn)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to create auth cache set command", log.ErrorField(err))

		return usecase.APIAuthRefreshOutput{}, err
	}

	if err := itr.authCache.Tx(ctx, []cache.TxSetCmd{aCmd}, []cache.TxDelCmd{cCmd}); err != nil {
		log.GetLogCtx(ctx).Warn("failed to transaction cache", log.ErrorField(err))

		return usecase.APIAuthRefreshOutput{}, err
	}

	return usecase.APIAuthRefreshOutput{
		AuthToken: at.ToToken(session.SessionID.ToSecret()),
	}, nil
}

func (itr *APIAuth) ChangePassword(
	ctx context.Context,
	input usecase.APIAuthChangePasswordInput,
) (usecase.APIAuthChangePasswordOutput, error) {
	if _, err := itr.authRPC.SignIn(ctx, input.EMail, input.OldPassword); err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign in with old password", log.ErrorField(err))

		return usecase.APIAuthChangePasswordOutput{}, errors.NewUnauthorizedError("failed to sign in", err)
	}

	if err := itr.authRPC.ChangePassword(ctx, input.UserID, input.NewPassword); err != nil {
		log.GetLogCtx(ctx).Warn("failed to change password", log.ErrorField(err))

		return usecase.APIAuthChangePasswordOutput{}, errors.NewUnknownError("failed to change password", err)
	}

	if _, err := itr.SignOutAll(ctx, usecase.APIAuthSignOutAllInput{
		UserID: input.UserID,
	}); err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign out all", log.ErrorField(err))

		return usecase.APIAuthChangePasswordOutput{}, errors.NewUnknownError("failed to sign out all", err)
	}

	output, err := itr.SignIn(ctx, usecase.APIAuthSignInInput{
		Secret:    input.Secret,
		EMail:     input.EMail,
		Password:  input.NewPassword,
		PublicKey: input.PublicKey,
		ExpiresIn: input.ExpiresIn,
	})
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign in with new password", log.ErrorField(err))

		return usecase.APIAuthChangePasswordOutput{}, errors.NewUnauthorizedError("failed to sign in", err)
	}

	return usecase.APIAuthChangePasswordOutput(output), nil
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
