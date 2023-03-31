package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	derr "github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
	"github.com/morning-night-guild/platform-app/pkg/trace"
)

var _ openapi.ServerInterface = (*Handler)(nil)

type Handler struct {
	key     string
	secret  auth.Secret
	auth    *Auth
	article *Article
	health  *Health
}

func New(
	key string,
	secret auth.Secret,
	auth *Auth,
	article *Article,
	health *Health,
) *Handler {
	return &Handler{
		key:     key,
		secret:  secret,
		auth:    auth,
		article: article,
		health:  health,
	}
}

type Auth struct {
	signUp       port.APIAuthSignUp
	signIn       port.APIAuthSignIn
	signOut      port.APIAuthSignOut
	verify       port.APIAuthVerify
	refresh      port.APIAuthRefresh
	generateCode port.APIAuthGenerateCode
	cookie       Cookie
}

type Cookie interface {
	Domain() string
	Secure() bool
	HTTPOnly() bool
	SameSite() http.SameSite
}

func NewAuth(
	signUp port.APIAuthSignUp,
	signIn port.APIAuthSignIn,
	signOut port.APIAuthSignOut,
	verify port.APIAuthVerify,
	refresh port.APIAuthRefresh,
	generateCode port.APIAuthGenerateCode,
	cookie Cookie,
) *Auth {
	return &Auth{
		signUp:       signUp,
		signIn:       signIn,
		signOut:      signOut,
		verify:       verify,
		refresh:      refresh,
		generateCode: generateCode,
		cookie:       cookie,
	}
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
