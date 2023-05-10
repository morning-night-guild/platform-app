package middleware

import (
	"net/http"
	"strconv"
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

		rw := newResponseWriter(w)

		next.ServeHTTP(rw, r.WithContext(ctx))

		logger.Info(
			"access-log",
			zap.String("method", r.Method),
			zap.String("path", r.RequestURI),
			zap.String("addr", r.RemoteAddr),
			zap.String("user-agent", r.Header["User-Agent"][0]),
			zap.String("status-code", strconv.Itoa(rw.StatusCode)),
			zap.String("elapsed", time.Since(now).String()),
			zap.Int64("elapsed(ms)", time.Since(now).Milliseconds()),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
