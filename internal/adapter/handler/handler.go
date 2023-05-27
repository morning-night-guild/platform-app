package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	derr "github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
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

func (hdl *Handler) HandleConnectError(ctx context.Context, err error) int { //nolint:cyclop
	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
		code := connect.CodeOf(connectErr)
		switch code {
		case connect.CodeInvalidArgument, connect.CodeOutOfRange, connect.CodeFailedPrecondition, connect.CodeAborted:
			log.GetLogCtx(ctx).Warn("invalid argument", log.ErrorField(err))

			return http.StatusBadRequest
		case connect.CodeUnauthenticated:
			log.GetLogCtx(ctx).Warn("unauthenticated", log.ErrorField(err))

			return http.StatusUnauthorized
		case connect.CodePermissionDenied:
			log.GetLogCtx(ctx).Warn("permission denied", log.ErrorField(err))

			return http.StatusForbidden
		case connect.CodeCanceled:
			log.GetLogCtx(ctx).Warn("canceled", log.ErrorField(err))

			return http.StatusRequestTimeout
		case connect.CodeAlreadyExists:
			log.GetLogCtx(ctx).Warn("already exists", log.ErrorField(err))

			return http.StatusConflict
		case connect.CodeNotFound:
			log.GetLogCtx(ctx).Warn("not found", log.ErrorField(err))

			return http.StatusNotFound
		case connect.CodeResourceExhausted:
			log.GetLogCtx(ctx).Warn("resource exhausted", log.ErrorField(err))

			return http.StatusTooManyRequests
		case connect.CodeUnimplemented:
			log.GetLogCtx(ctx).Error("unimplemented", log.ErrorField(err))

			return http.StatusNotImplemented
		case connect.CodeUnknown, connect.CodeInternal, connect.CodeDataLoss:
			log.GetLogCtx(ctx).Error("unknown", log.ErrorField(err))

			return http.StatusInternalServerError
		case connect.CodeUnavailable:
			log.GetLogCtx(ctx).Error("unavailable", log.ErrorField(err))

			return http.StatusServiceUnavailable
		case connect.CodeDeadlineExceeded:
			log.GetLogCtx(ctx).Error("deadline exceeded", log.ErrorField(err))

			return http.StatusGatewayTimeout
		}
	}

	log.GetLogCtx(ctx).Error("unknown error", log.ErrorField(err))

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

func NewRequest[T any](ctx context.Context, msg *T) *connect.Request[T] {
	req := connect.NewRequest(msg)

	req.Header().Set("tid", trace.GetTIDCtx(ctx))

	return req
}

type Tokens struct {
	SessionToken auth.SessionToken
	AuthToken    auth.AuthToken
}

func (hdl *Handler) ExtractTokens(
	ctx context.Context,
	r *http.Request,
) (Tokens, error) {
	sessionTokenCookie, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get session token cookie", log.ErrorField(err))

		return Tokens{}, derr.NewUnauthorizedError("failed to get session token cookie", err)
	}

	sessionToken, err := auth.ParseSessionToken(sessionTokenCookie.Value, hdl.secret)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new session token", log.ErrorField(err))

		return Tokens{}, derr.NewUnauthorizedError("failed to new session token", err)
	}

	authTokenCookie, err := r.Cookie(auth.AuthTokenKey)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to get auth token cookie", log.ErrorField(err))

		return Tokens{}, derr.NewUnauthorizedError("failed to get auth token cookie", err)
	}

	authToken, err := auth.ParseAuthToken(authTokenCookie.Value, sessionToken.ToSecret(hdl.secret))
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to new auth token", log.ErrorField(err))

		return Tokens{}, derr.NewUnauthorizedError("failed to new auth token", err)
	}

	return Tokens{
		AuthToken:    authToken,
		SessionToken: sessionToken,
	}, nil
}

func (hdl *Handler) ExtractUserID(
	ctx context.Context,
	r *http.Request,
) (user.ID, error) {
	tokens, err := hdl.ExtractTokens(ctx, r)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to extract tokens", log.ErrorField(err))

		return user.GenerateZeroID(), err
	}

	uid := tokens.AuthToken.UserID()

	return uid, nil
}
