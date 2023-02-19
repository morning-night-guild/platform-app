package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

var _ openapi.ServerInterface = (*API)(nil)

type API struct {
	key     string
	connect *Connect
}

func New(
	key string,
	connect *Connect,
) *API {
	return &API{
		key:     key,
		connect: connect,
	}
}

func (api *API) HandleConnectError(ctx context.Context, err error) int {
	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
		code := connect.CodeOf(connectErr)
		if code == connect.CodeInvalidArgument {
			log.GetLogCtx(ctx).Warn("invalid argument.", log.ErrorField(err))

			return http.StatusBadRequest
		}
	}

	log.GetLogCtx(ctx).Error("failed to share article", log.ErrorField(err))

	return http.StatusInternalServerError
}
