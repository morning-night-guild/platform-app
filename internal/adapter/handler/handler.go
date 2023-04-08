package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	derr "github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
	"github.com/morning-night-guild/platform-app/pkg/trace"
)

var _ openapi.ServerInterface = (*Handler)(nil)

type Handler struct {
	key     string
	secret  auth.Secret
	cookie  Cookie
	auth    usecase.APIAuth
	article usecase.APIArticle
	health  usecase.APIHealth
}

func New(
	key string,
	secret auth.Secret,
	cookie Cookie,
	auth usecase.APIAuth,
	article usecase.APIArticle,
	health usecase.APIHealth,
) *Handler {
	return &Handler{
		key:     key,
		secret:  secret,
		cookie:  cookie,
		auth:    auth,
		article: article,
		health:  health,
	}
}

func (hdl *Handler) HandleConnectError(ctx context.Context, err error) int {
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

func (hdl *Handler) HandleErrorStatus(
	w http.ResponseWriter,
	err error,
) {
	switch {
	case derr.AsValidationError(err):
		w.WriteHeader(http.StatusBadRequest)
	case derr.AsURLError(err):
		w.WriteHeader(http.StatusBadRequest)
	case derr.AsUnauthorizedError(err):
		w.WriteHeader(http.StatusUnauthorized)
	case derr.AsNotFoundError(err):
		w.WriteHeader(http.StatusNotFound)
	case derr.AsUnknownError(err):
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (hdl *Handler) PointerToString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func (hdl *Handler) StringToPointer(s string) *string {
	return &s
}

func NewRequestWithTID[T any](ctx context.Context, msg *T) *connect.Request[T] {
	req := connect.NewRequest(msg)

	req.Header().Set("tid", trace.GetTIDCtx(ctx))

	return req
}
