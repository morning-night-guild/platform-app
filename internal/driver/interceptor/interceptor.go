package interceptor

import (
	"context"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

// New.
func New() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			now := time.Now()

			tid := req.Header().Get("tid")

			ctx = trace.SetTIDCtx(ctx, tid)

			ctx = log.SetLogCtx(ctx, trace.GetTIDCtx(ctx))

			logger := log.GetLogCtx(ctx)

			res, err := next(ctx, req)

			logger.Info(
				"access log",
				zap.String("path", req.Spec().Procedure),
				zap.String("protocol", req.Peer().Protocol),
				zap.String("addr", req.Peer().Addr),
				zap.String("ua", req.Header().Get("User-Agent")),
				zap.String("code", status.Code(err).String()),
				zap.String("elapsed", time.Since(now).String()),
				zap.Int64("elapsed(ns)", time.Since(now).Nanoseconds()),
			)

			return res, err
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
