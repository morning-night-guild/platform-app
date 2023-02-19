package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/morning-night-guild/platform-app/pkg/log"
)

const (
	shutdownTime      = 10
	readHeaderTimeout = 30 * time.Second
)

// Server.
type Server struct {
	*http.Server
}

// NewServer
// 引数nrはnilでも動作可能（NewRelicへレポートが送信されない）.
func NewServer(
	port string,
	handler http.Handler,
) *Server {
	s := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &Server{
		Server: s,
	}
}

// Run.
func (srv *Server) Run() {
	log.Log().Sugar().Infof("server running on %s", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Log().Sugar().Errorf("server closed with error: %s", err.Error())

			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	log.Log().Sugar().Infof("SIGNAL %d received, then shutting down...\n", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Log().Sugar().Infof("failed to gracefully shutdown:", err)
	}

	log.Log().Info("server shutdown")
}
