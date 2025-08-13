package routes

import (
	"github.com/NERFTHISPLS/rest-todo-list/internal/handlers/tasks"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(app *fiber.App, dbpool *pgxpool.Pool) {
	taskHandler := tasks.NewHandler(dbpool)

	app.Get("/tasks", taskHandler.List)
	app.Post("/tasks", taskHandler.Create)
	app.Put("/tasks/:id", taskHandler.Update)
	app.Delete("/tasks/:id", taskHandler.Delete)
}
