package handler

import (
	"net/http"

	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

// apiヘルスチェック
// (GET /v1/health/api).
func (hdl *Handler) V1HealthAPI(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// coreヘルスチェック
// (GET /v1/health/core).
func (hdl *Handler) V1HealthCore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := usecase.APIHealthCheckInput{}

	if _, err := hdl.health.Check(ctx, input); err != nil {
		log.GetLogCtx(ctx).Error("failed to check health core", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
