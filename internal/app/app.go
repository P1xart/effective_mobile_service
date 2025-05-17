package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/P1xart/effective_mobile_service/internal/config"
	v1 "github.com/P1xart/effective_mobile_service/internal/controller/v1"
	"github.com/P1xart/effective_mobile_service/internal/repo"
	"github.com/P1xart/effective_mobile_service/internal/service"
	"github.com/P1xart/effective_mobile_service/pkg/httpsrv"
	"github.com/P1xart/effective_mobile_service/pkg/logger"
	"github.com/P1xart/effective_mobile_service/pkg/postgresql"

	"github.com/gin-gonic/gin"
)

func Run() {
	log := logger.New()
	log.Debug("app starting")

	cfg, err := config.New(log)
	if err != nil {
		log.Error("failed to init config", logger.Error(err))
		os.Exit(1)
	}

	log.Debug("postgresql starting")
	postgres, err := postgresql.New(log, &cfg.Postgresql)
	if err != nil {
		log.Error("failed to init postgresql", logger.Error(err))
		os.Exit(1)
	}

	log.Debug("repositories init")
	repositories := repo.NewRepositories(log, postgres)

	log.Debug("services init")
	services := service.NewServices(&service.Dependencies{
		Log:      log,
		Repos:    repositories,
	}, &cfg.API)

	r := router()
	v1.NewRouter(log, r, services)

	log.Debug("server starting")
	server := httpsrv.New(log, cfg, r)

	log.Info("configuring graceful shutdown")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("application got signal", slog.String("signal", s.String()))
	case err = <-server.Notify():
		log.Error("http server error", logger.Error(err))
	}

	if err = server.Shutdown(); err != nil {
		log.Error("failed to shutdown http server", logger.Error(err))
	}

	postgres.Close()
}

func router() *gin.Engine {
	var r *gin.Engine

	if env := os.Getenv("APP_ENV"); env == "prod" {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(gin.Recovery())
	} else {
		r = gin.Default()
	}

	return r
}
