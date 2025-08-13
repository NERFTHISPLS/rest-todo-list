package server

import (
	"fmt"
	"log"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
	"github.com/NERFTHISPLS/rest-todo-list/internal/repository"
	"github.com/NERFTHISPLS/rest-todo-list/internal/server/routes"
	"github.com/gofiber/fiber/v2"
)

func Setup(cfg *config.ConfServer, repo *repository.TaskRepository) error {
	log.Printf("starting server on :%d", cfg.Port)

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.TimeoutRead,
		WriteTimeout: cfg.TimeoutWrite,
		IdleTimeout:  cfg.TimeoutIdle,
	})
	serverPort := fmt.Sprintf(":%d", cfg.Port)

	routes.Setup(app, repo)

	return app.Listen(serverPort)
}
