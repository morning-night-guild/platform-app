package api

import (
	"net/http"

	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

func (api *API) V1HealthAPI(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func (api *API) V1HealthCore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := port.APIHealthCheckInput{}

	if _, err := api.health.check.Execute(ctx, input); err != nil {
		log.GetLogCtx(ctx).Error("failed to check health core", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
