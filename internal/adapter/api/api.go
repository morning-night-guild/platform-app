package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
	"github.com/morning-night-guild/platform-app/pkg/trace"
)

var _ openapi.ServerInterface = (*API)(nil)

type API struct {
	key     string
	article *Article
	health  *Health
}

type Article struct {
	list  port.APIArticleList
	share port.APIArticleShare
}

func NewArticle(
	list port.APIArticleList,
	share port.APIArticleShare,
) *Article {
	return &Article{
		list:  list,
		share: share,
	}
}

type Health struct {
	check port.APIHealthCheck
}

func NewHealth(
	check port.APIHealthCheck,
) *Health {
	return &Health{
		check: check,
	}
}

func New(
	key string,
	article *Article,
	health *Health,
) *API {
	return &API{
		key:     key,
		article: article,
		health:  health,
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

func (api *API) PointerToString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func NewRequestWithTID[T any](ctx context.Context, msg *T) *connect.Request[T] {
	req := connect.NewRequest(msg)

	req.Header().Set("tid", trace.GetTIDCtx(ctx))

	return req
}
