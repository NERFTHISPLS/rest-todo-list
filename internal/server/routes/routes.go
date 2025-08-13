package routes

import (
	"github.com/NERFTHISPLS/rest-todo-list/internal/handlers/tasks"
	"github.com/NERFTHISPLS/rest-todo-list/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, repo *repository.TaskRepository) {
	taskHandler := tasks.NewHandler(repo)

	app.Get("/tasks", taskHandler.List)
	app.Post("/tasks", taskHandler.Create)
	app.Put("/tasks/:id", taskHandler.Update)
	app.Delete("/tasks/:id", taskHandler.Delete)
}
