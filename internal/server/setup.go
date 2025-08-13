package server

import (
	"fmt"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
	"github.com/gofiber/fiber/v2"
)

func Setup(cfg *config.ConfServer) error {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.TimeoutRead,
		WriteTimeout: cfg.TimeoutWrite,
		IdleTimeout:  cfg.TimeoutIdle,
	})
	serverPort := fmt.Sprintf(":%d", cfg.Port)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!!!")
	})

	return app.Listen(serverPort)
}
