package server

import (
	"fmt"
	"log"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
	"github.com/gofiber/fiber/v2"
)

func Setup(cfg *config.ConfServer) error {
	log.Printf("starting server on :%d", cfg.Port)

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
