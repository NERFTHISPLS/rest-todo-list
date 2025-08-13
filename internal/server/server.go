package server

import (
	"fmt"
	"log"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
	"github.com/NERFTHISPLS/rest-todo-list/internal/server/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(cfg *config.ConfServer, dbpool *pgxpool.Pool) error {
	log.Printf("starting server on :%d", cfg.Port)

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.TimeoutRead,
		WriteTimeout: cfg.TimeoutWrite,
		IdleTimeout:  cfg.TimeoutIdle,
	})
	serverPort := fmt.Sprintf(":%d", cfg.Port)

	routes.Setup(app, dbpool)

	return app.Listen(serverPort)
}
