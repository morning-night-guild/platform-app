package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

const path = "/"

// リフレッシュ
// (GET /v1/auth/refresh).
func (hdl *Handler) V1AuthRefresh(w http.ResponseWriter, r *http.Request, params openapi.V1AuthRefreshParams) {
	ctx := r.Context()

	codeID, err := auth.NewCodeID(params.Code)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new code id", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	signature, err := auth.NewSignature(params.Signature)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new signature", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	sessionTokenCookie, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get session token", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	sessionToken, err := auth.ParseSessionToken(sessionTokenCookie.Value, hdl.secret)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new session token", log.ErrorField(err))

		hdl.HandleErrorStatus(w, err)

		return
	}

	expiresIn := auth.DefaultExpiresIn
	if params.ExpiresIn != nil {
		expiresIn = auth.ExpiresIn(*params.ExpiresIn)
	}

	input := usecase.APIAuthRefreshInput{
		CodeID:    codeID,
		Signature: signature,
		SessionID: sessionToken.ID(hdl.secret),
		ExpiresIn: expiresIn,
	}

	output, err := hdl.auth.Refresh(ctx, input)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to refresh", log.ErrorField(err))

		hdl.HandleErrorStatus(w, err)

		return
	}

	expires := model.DefaultAuthExpiresIn
	if params.ExpiresIn != nil {
		expires = time.Duration(*params.ExpiresIn * int(time.Second))
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthTokenKey,
		Value:    output.AuthToken.String(),
		Path:     path,
		Domain:   hdl.cookie.Domain(),
		Expires:  time.Now().Add(expires),
		Secure:   hdl.cookie.Secure(),
		HttpOnly: true,
		SameSite: hdl.cookie.SameSite(),
	})
}

// サインイン
// (POST /v1/auth/signin).
func (hdl *Handler) V1AuthSignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req openapi.V1AuthSignInRequestSchema

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.GetLogCtx(ctx).Warn("failed to decode request body", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	key, err := auth.DecodePublicKey(req.PublicKey)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to decode public key", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	email, err := auth.NewEMail(string(req.Email))
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new email", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	password, err := auth.NewPassword(req.Password)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new password", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	expiresIn := auth.DefaultExpiresIn
	if req.ExpiresIn != nil {
		expiresIn, err = auth.NewExpiresIn(*req.ExpiresIn)
		if err != nil {
			log.GetLogCtx(ctx).Warn("failed to new expires in", log.ErrorField(err))

			w.WriteHeader(http.StatusBadRequest)

			return
		}
	}

	input := usecase.APIAuthSignInInput{
		Secret:    hdl.secret,
		PublicKey: key,
		EMail:     email,
		Password:  password,
		ExpiresIn: expiresIn,
	}

	output, err := hdl.auth.SignIn(ctx, input)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign in", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthTokenKey,
		Value:    output.AuthToken.String(),
		Path:     path,
		Domain:   hdl.cookie.Domain(),
		Expires:  output.Auth.ExpiresAt,
		Secure:   hdl.cookie.Secure(),
		HttpOnly: true,
		SameSite: hdl.cookie.SameSite(),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionTokenKey,
		Value:    output.SessionToken.String(),
		Path:     path,
		Domain:   hdl.cookie.Domain(),
		Expires:  output.Auth.IssuedAt.Add(model.DefaultSessionExpiresIn),
		Secure:   hdl.cookie.Secure(),
		HttpOnly: true,
		SameSite: hdl.cookie.SameSite(),
	})
}

// サインアウト
// (GET /v1/auth/signout).
func (hdl *Handler) V1AuthSignOut(_ http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokens, err := hdl.ExtractTokens(ctx, r)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to extract tokens", log.ErrorField(err))

		return
	}

	input := usecase.APIAuthSignOutInput{
		UserID:    tokens.AuthToken.UserID(),
		SessionID: tokens.SessionToken.ID(hdl.secret),
	}

	if _, err := hdl.auth.SignOut(ctx, input); err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign out", log.ErrorField(err))

		return
	}
}

// サインアウトオール
// (GET /v1/auth/signout/all).
func (hdl *Handler) V1AuthSignOutAll(_ http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokens, err := hdl.ExtractTokens(ctx, r)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to extract tokens", log.ErrorField(err))

		return
	}

	input := usecase.APIAuthSignOutAllInput{
		UserID: tokens.AuthToken.UserID(),
	}

	if _, err := hdl.auth.SignOutAll(ctx, input); err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign out", log.ErrorField(err))

		return
	}
}

// サインアップ
// (POST /v1/auth/signup).
func (hdl *Handler) V1AuthSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key := r.Header.Get("Api-Key")
	if key != hdl.key {
		log.GetLogCtx(ctx).Warn(fmt.Sprintf("invalid api key. api key = %s", key))

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	var body openapi.V1AuthSignUpRequestSchema

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.GetLogCtx(ctx).Warn("failed to decode request body", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	email, err := auth.NewEMail(string(body.Email))
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new email", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	password, err := auth.NewPassword(body.Password)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new password", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	input := usecase.APIAuthSignUpInput{
		EMail:    email,
		Password: password,
	}

	if _, err := hdl.auth.SignUp(ctx, input); err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign up", log.ErrorField(err))

		hdl.HandleErrorStatus(w, err)

		return
	}
}

// 検証
// (GET /v1/auth/verify).
func (hdl *Handler) V1AuthVerify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 取り出すトークンの種類によってエラーレスポンスが異なるので
	// ExtractTokens() は使わない

	sessionTokenCookie, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get session token cookie", log.ErrorField(err))

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	sessionToken, err := auth.ParseSessionToken(sessionTokenCookie.Value, hdl.secret)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new session token", log.ErrorField(err))

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	authTokenCookie, err := r.Cookie(auth.AuthTokenKey)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get auth token cookie", log.ErrorField(err))

		hdl.unauthorize(ctx, w, sessionToken)

		return
	}

	authToken, err := auth.ParseAuthToken(authTokenCookie.Value, sessionToken.ToSecret(hdl.secret))
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new auth token", log.ErrorField(err))

		hdl.unauthorize(ctx, w, sessionToken)

		return
	}

	input := usecase.APIAuthVerifyInput{
		UserID:    authToken.UserID(),
		SessionID: sessionToken.ID(hdl.secret),
	}

	if _, err := hdl.auth.Verify(ctx, input); err != nil {
		log.GetLogCtx(ctx).Warn("failed to verify", log.ErrorField(err))

		hdl.unauthorize(ctx, w, sessionToken)

		return
	}

	_, _ = w.Write([]byte("OK"))
}

func (hdl *Handler) unauthorize(
	ctx context.Context,
	w http.ResponseWriter,
	sessionToken auth.SessionToken,
) {
	input := usecase.APIAuthGenerateCodeInput{
		SessionID: sessionToken.ID(hdl.secret),
	}

	output, err := hdl.auth.GenerateCode(ctx, input)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to generate code", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	res := openapi.V1AuthVerifyUnauthorizedResponseSchema{
		Code: output.Code.CodeID.Value(),
	}

	w.WriteHeader(http.StatusUnauthorized)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.GetLogCtx(ctx).Warn("failed to encode response body", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

// パスワード変更
// (PUT /v1/auth/password).
func (hdl *Handler) V1AuthChangePassword( //nolint:cyclop
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	tokens, err := hdl.ExtractTokens(ctx, r)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to extract tokens", log.ErrorField(err))

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	var req openapi.V1AuthChangePasswordRequestSchema

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.GetLogCtx(ctx).Warn("failed to decode request body", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	key, err := auth.DecodePublicKey(req.PublicKey)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to decode public key", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	email, err := auth.NewEMail(string(req.Email))
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new email", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	oldPassword, err := auth.NewPassword(req.OldPassword)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new old password", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	newPassword, err := auth.NewPassword(req.NewPassword)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new password", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if oldPassword.Equal(newPassword) {
		log.GetLogCtx(ctx).Warn("old password and new password are equal")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	expiresIn := auth.DefaultExpiresIn
	if req.ExpiresIn != nil {
		expiresIn, err = auth.NewExpiresIn(*req.ExpiresIn)
		if err != nil {
			log.GetLogCtx(ctx).Warn("failed to new expires in", log.ErrorField(err))

			w.WriteHeader(http.StatusBadRequest)

			return
		}
	}

	input := usecase.APIAuthChangePasswordInput{
		UserID:      tokens.AuthToken.UserID(),
		Secret:      hdl.secret,
		PublicKey:   key,
		ExpiresIn:   expiresIn,
		EMail:       email,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	output, err := hdl.auth.ChangePassword(ctx, input)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to change password", log.ErrorField(err))

		hdl.HandleErrorStatus(w, err)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthTokenKey,
		Value:    output.AuthToken.String(),
		Path:     path,
		Domain:   hdl.cookie.Domain(),
		Expires:  output.Auth.ExpiresAt,
		Secure:   hdl.cookie.Secure(),
		HttpOnly: true,
		SameSite: hdl.cookie.SameSite(),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionTokenKey,
		Value:    output.SessionToken.String(),
		Path:     path,
		Domain:   hdl.cookie.Domain(),
		Expires:  output.Auth.IssuedAt.Add(model.DefaultSessionExpiresIn),
		Secure:   hdl.cookie.Secure(),
		HttpOnly: true,
		SameSite: hdl.cookie.SameSite(),
	})
}
