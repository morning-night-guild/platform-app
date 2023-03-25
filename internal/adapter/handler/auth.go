package handler

import (
	"net/http"

	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

// リフレッシュ
// (GET /v1/auth/refresh)
func (hand *Handler) V1AuthRefresh(w http.ResponseWriter, r *http.Request, params openapi.V1AuthRefreshParams) {
}

// サインイン
// (POST /v1/auth/signin)
func (hand *Handler) V1AuthSignIn(w http.ResponseWriter, r *http.Request) {
}

// サインアウト
// (GET /v1/auth/signout)
func (hand *Handler) V1AuthSignOut(w http.ResponseWriter, r *http.Request) {
}

// サインアップ
// (POST /v1/auth/signup)
func (hand *Handler) V1AuthSignUp(w http.ResponseWriter, r *http.Request) {
}

// 検証
// (GET /v1/auth/verify)
func (hand *Handler) V1AuthVerify(w http.ResponseWriter, r *http.Request) {
}
