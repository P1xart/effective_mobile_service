package logger

import (
	"log/slog"
	"os"
	"strings"
)

const defaultLevel = slog.LevelDebug

func New() *slog.Logger {
	var handler slog.Handler

	if env := os.Getenv("APP_ENV"); env == "prod" {
		handler = prodHandler()
	} else if env == "dev" {
		handler = devHandler()
	} else {
		return slog.Default()
	}

	return slog.New(handler)
}

func prodHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: configLevel(),
	})
}

func devHandler() slog.Handler {
	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     configLevel(),
		AddSource: true,
	})
}

func configLevel() slog.Level {
	var logLevel slog.Level

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = defaultLevel
	}

	return logLevel
}
