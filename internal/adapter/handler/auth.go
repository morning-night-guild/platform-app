package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

const path = "/"

// リフレッシュ
// (GET /v1/auth/refresh).
//
//nolint:funlen
func (hdl *Handler) V1AuthRefresh(w http.ResponseWriter, r *http.Request, params openapi.V1AuthRefreshParams) {
	ctx := r.Context()

	codeID, err := auth.NewCodeID(params.Code)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new code id", log.ErrorField(err))

		hdl.HandleErrorStatus(w, err)

		return
	}

	signature, err := auth.NewSignature(params.Signature)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new signature", log.ErrorField(err))

		hdl.HandleErrorStatus(w, err)

		return
	}

	sessionTokenCookie, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get session token", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	sessionToken, err := auth.NewSessionToken(sessionTokenCookie.Value, hdl.secret)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new session token", log.ErrorField(err))

		hdl.HandleErrorStatus(w, err)

		return
	}

	input := port.APIAuthRefreshInput{
		CodeID:       codeID,
		Signature:    signature,
		SessionToken: sessionToken,
	}

	output, err := hdl.auth.refresh.Execute(ctx, input)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to refresh", log.ErrorField(err))

		hdl.HandleErrorStatus(w, err)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthTokenKey,
		Value:    output.AuthToken.String(),
		Path:     path,
		Domain:   hdl.auth.cookie.Domain(),
		Expires:  time.Now().Add(model.DefaultAuthExpiresIn),
		Secure:   hdl.auth.cookie.Secure(),
		HttpOnly: hdl.auth.cookie.HTTPOnly(),
		SameSite: hdl.auth.cookie.SameSite(),
	})
}

// サインイン
// (POST /v1/auth/signin).
//
//nolint:funlen
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

		hdl.HandleErrorStatus(w, err)

		return
	}

	email, err := auth.NewEMail(string(req.Email))
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new email", log.ErrorField(err))

		hdl.HandleErrorStatus(w, err)

		return
	}

	password, err := auth.NewPassword(req.Password)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new password", log.ErrorField(err))

		hdl.HandleErrorStatus(w, err)

		return
	}

	expiresIn := auth.DefaultExpiresIn
	if req.ExpiresIn != nil {
		expiresIn, err = auth.NewExpiresIn(*req.ExpiresIn)
		if err != nil {
			log.GetLogCtx(ctx).Warn("failed to new expires in", log.ErrorField(err))

			hdl.HandleErrorStatus(w, err)

			return
		}
	}

	input := port.APIAuthSignInInput{
		PublicKey: key,
		EMail:     email,
		Password:  password,
		ExpiresIn: expiresIn,
	}

	output, err := hdl.auth.signIn.Execute(ctx, input)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign in", log.ErrorField(err))

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	now := time.Now()

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthTokenKey,
		Value:    output.AuthToken.String(),
		Path:     path,
		Domain:   hdl.auth.cookie.Domain(),
		Expires:  now.Add(model.DefaultAuthExpiresIn),
		Secure:   hdl.auth.cookie.Secure(),
		HttpOnly: hdl.auth.cookie.HTTPOnly(),
		SameSite: hdl.auth.cookie.SameSite(),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionTokenKey,
		Value:    output.SessionToken.String(),
		Path:     path,
		Domain:   hdl.auth.cookie.Domain(),
		Expires:  now.Add(model.DefaultSessionExpiresIn),
		Secure:   hdl.auth.cookie.Secure(),
		HttpOnly: hdl.auth.cookie.HTTPOnly(),
		SameSite: hdl.auth.cookie.SameSite(),
	})
}

// サインアウト
// (GET /v1/auth/signout).
func (hdl *Handler) V1AuthSignOut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionTokenCookie, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get session token cookie", log.ErrorField(err))

		return
	}

	sessionToken, err := auth.NewSessionToken(sessionTokenCookie.Value, hdl.secret)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new session token", log.ErrorField(err))

		return
	}

	authTokenCookie, err := r.Cookie(auth.AuthTokenKey)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get auth token cookie", log.ErrorField(err))

		return
	}

	authToken, err := auth.NewAuthToken(authTokenCookie.Value, sessionToken.ToSecret(hdl.secret))
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new auth token", log.ErrorField(err))

		return
	}

	input := port.APIAuthSignOutInput{
		AuthToken:    authToken,
		SessionToken: sessionToken,
	}

	if _, err := hdl.auth.signOut.Execute(ctx, input); err != nil {
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

		hdl.HandleErrorStatus(w, err)

		return
	}

	password, err := auth.NewPassword(body.Password)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new password", log.ErrorField(err))

		hdl.HandleErrorStatus(w, err)

		return
	}

	input := port.APIAuthSignUpInput{
		EMail:    email,
		Password: password,
	}

	if _, err := hdl.auth.signUp.Execute(ctx, input); err != nil {
		log.GetLogCtx(ctx).Warn("failed to sign up", log.ErrorField(err))

		hdl.HandleErrorStatus(w, err)

		return
	}
}

// 検証
// (GET /v1/auth/verify).
func (hdl *Handler) V1AuthVerify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionTokenCookie, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get session token cookie", log.ErrorField(err))

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	sessionToken, err := auth.NewSessionToken(sessionTokenCookie.Value, hdl.secret)
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

	authToken, err := auth.NewAuthToken(authTokenCookie.Value, sessionToken.ToSecret(hdl.secret))
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new auth token", log.ErrorField(err))

		hdl.unauthorize(ctx, w, sessionToken)

		return
	}

	input := port.APIAuthVerifyInput{
		AuthToken:    authToken,
		SessionToken: sessionToken,
	}

	if _, err := hdl.auth.verify.Execute(ctx, input); err != nil {
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
	input := port.APIAuthGenerateCodeInput{
		SessionToken: sessionToken,
	}

	output, err := hdl.auth.generateCode.Execute(ctx, input)
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
