package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
	"github.com/NERFTHISPLS/rest-todo-list/internal/repository"
	"github.com/NERFTHISPLS/rest-todo-list/internal/server/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Setup(cfg *config.ConfServer, repo *repository.TaskRepository) error {
	slog.Info("starting server", "port", cfg.Port)

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.TimeoutRead,
		WriteTimeout: cfg.TimeoutWrite,
		IdleTimeout:  cfg.TimeoutIdle,
	})

	logFormat := "[${time}] ${status} - ${latency} ${method} ${path} - ${ip}\n"
	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		logFormat = "[${time}] ${status} - ${latency} ${method} ${path} - ${ip} - ${user_agent}\n"
	}

	app.Use(logger.New(logger.Config{
		Format:     logFormat,
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Local",
		Next: func(c *fiber.Ctx) bool {
			if !slog.Default().Enabled(context.Background(), slog.LevelInfo) {
				return c.Response().StatusCode() >= 400
			}

			return false
		},
	}))

	serverPort := fmt.Sprintf(":%d", cfg.Port)

	routes.Setup(app, repo)

	slog.Info("server configured successfully", "port", cfg.Port)
	return app.Listen(serverPort)
}
