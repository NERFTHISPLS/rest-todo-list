package routes

import (
	"github.com/NERFTHISPLS/rest-todo-list/internal/handlers/tasks"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/tasks", tasks.List)
	app.Post("/tasks", tasks.Create)
	app.Put("/tasks/:id", tasks.Update)
	app.Delete("/tasks/:id", tasks.Delete)
}
