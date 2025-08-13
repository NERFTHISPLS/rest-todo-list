package tasks

import (
	"log"

	"github.com/NERFTHISPLS/rest-todo-list/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	dbPool *pgxpool.Pool
}

func NewHandler(dbpool *pgxpool.Pool) *Handler {
	return &Handler{
		dbPool: dbpool,
	}
}

func (h *Handler) List(c *fiber.Ctx) error {
	if h.dbPool == nil {
		return c.Status(500).SendString("database connection is not initialized")
	}

	rows, err := h.dbPool.Query(c.Context(), "SELECT * FROM tasks")
	if err != nil {
		log.Printf("error while running the query on the database: %s\n", err)
		return c.Status(500).SendString(err.Error())
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task

		if err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			log.Printf("error while scanning the rows: %s", err)
			return c.Status(500).SendString(err.Error())
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(tasks)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) Update(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	return nil
}
