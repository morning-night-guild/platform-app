package interceptor

import (
	"context"
	"errors"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/internal/adapter/external"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/trace"
	"go.uber.org/zap"
)

// New.
func New() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			now := time.Now()

			tid := req.Header().Get(external.HeaderTID)

			ctx = trace.SetTIDCtx(ctx, tid)

			ctx = log.SetLogCtx(ctx, trace.GetTIDCtx(ctx))

			uid := req.Header().Get(external.HeaderUID)

			userID, _ := user.NewID(uid)

			ctx = user.SetUIDCtx(ctx, userID)

			logger := log.GetLogCtx(ctx)

			res, err := next(ctx, req)

			code := func(err error) string {
				if err == nil {
					return "ok"
				}

				connectErr := new(connect.Error)

				if !errors.As(err, &connectErr) {
					return "unknown"
				}

				return connect.CodeOf(connectErr).String()
			}

			logger.Info(
				"access-log",
				zap.String("uid", uid),
				zap.String("path", req.Spec().Procedure),
				zap.String("protocol", req.Peer().Protocol),
				zap.String("addr", req.Peer().Addr),
				zap.String("user-agent", req.Header().Get("User-Agent")),
				zap.String("status-code", code(err)),
				zap.String("elapsed", time.Since(now).String()),
				zap.Int64("elapsed(ms)", time.Since(now).Milliseconds()),
			)

			return res, err
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
