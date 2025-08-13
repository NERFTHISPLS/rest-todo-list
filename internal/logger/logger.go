package logger

import (
	"log/slog"
	"os"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
)

func Setup(cfg *config.Conf) {
	var logLevel slog.Level

	switch cfg.Env {
	case "prod":
		logLevel = slog.LevelWarn
		slog.Info("running in production mode", "log_level", "warn")
	case "dev":
		logLevel = slog.LevelDebug
		slog.Info("running in development mode", "log_level", "debug")
	default:
		logLevel = slog.LevelInfo
		slog.Info("running in default mode", "log_level", "info")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	slog.SetDefault(logger)

	slog.Info("app start", "environment", cfg.Env, "log_level", logLevel)
}
