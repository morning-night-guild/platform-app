package middleware

import (
	"net/http"
	"time"

	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/trace"
	"go.uber.org/zap"
)

// Middleware.
type Middleware struct{}

// New.
func New() *Middleware {
	return &Middleware{}
}

// Handle.
func (middle *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()

		ctx := trace.WithTIDCtx(r.Context())

		ctx = log.SetLogCtx(ctx, trace.GetTIDCtx(ctx))

		logger := log.GetLogCtx(ctx)

		next.ServeHTTP(w, r.WithContext(ctx))

		logger.Info(
			"access log",
			zap.String("method", r.Method),
			zap.String("path", r.RequestURI),
			zap.String("addr", r.RemoteAddr),
			zap.String("ua", r.Header["User-Agent"][0]),
			// zap.String("code", status.Code(err).String()),
			zap.String("elapsed", time.Since(now).String()),
			zap.Int64("elapsed(ns)", time.Since(now).Nanoseconds()),
		)
	})
}
