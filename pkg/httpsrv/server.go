package httpsrv

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/P1xart/effective_mobile_service/internal/config"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	log    *slog.Logger
	cfg    *config.Config
	notify chan error

	base *http.Server
}

func New(log *slog.Logger, cfg *config.Config, router http.Handler) *Server {
	log = log.With(slog.String("component", "http server"))

	srv := &Server{
		log:    log,
		cfg:    cfg,
		notify: make(chan error),
		base: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
			Handler:      router,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
		},
	}

	srv.start()

	return srv
}

func (srv *Server) start() {
	go func() {
		srv.notify <- srv.base.ListenAndServe()
		close(srv.notify)
	}()
}

func (srv *Server) Notify() <-chan error {
	return srv.notify
}

func (srv *Server) Shutdown() error {
	srv.log.Debug("shutdown http server")
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	return srv.base.Shutdown(ctx)
}
