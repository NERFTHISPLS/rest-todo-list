package tasks

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/NERFTHISPLS/rest-todo-list/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
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
		return c.Status(fiber.StatusInternalServerError).SendString("database connection is not initialized")
	}

	query := `SELECT * FROM tasks`
	rows, err := h.dbPool.Query(c.Context(), query)
	if err != nil {
		log.Printf("error while running the query on the database: %s\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task

		if err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			log.Printf("error while scanning the rows: %s", err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(tasks)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	task := &models.Task{}
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if strings.TrimSpace(task.Title) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
	}

	if strings.TrimSpace(task.Status) == "" {
		task.Status = models.DefaultTaskStatus
	}

	query := `INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id`
	if err := h.dbPool.QueryRow(
		c.Context(), query, task.Title, task.Description, task.Status,
	).Scan(&task.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create task"})
	}

	return c.Status(fiber.StatusOK).JSON(task)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	task := &models.Task{}
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	setClauses := []string{}
	args := []any{}
	argPos := 1

	if task.Title != "" {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argPos))
		args = append(args, task.Title)
		argPos++
	}

	if task.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argPos))
		args = append(args, task.Description)
		argPos++
	}

	if task.Status != "" {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argPos))
		args = append(args, task.Status)
		argPos++
	}

	if len(setClauses) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no fields to update"})
	}

	setClauses = append(setClauses, "updated_at = now()")

	query := fmt.Sprintf(`
		UPDATE tasks
		SET %s
		WHERE id = $%d
		RETURNING id, title, description, status, created_at, updated_at
	`, strings.Join(setClauses, ", "), argPos)

	args = append(args, id)

	row := h.dbPool.QueryRow(c.Context(), query, args...)

	updatedTask := &models.Task{}

	if err := row.Scan(
		&updatedTask.ID,
		&updatedTask.Title,
		&updatedTask.Description,
		&updatedTask.Status,
		&updatedTask.CreatedAt,
		&updatedTask.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update task"})
	}

	return c.Status(fiber.StatusOK).JSON(updatedTask)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	query := `DELETE FROM tasks WHERE id = $1`
	cmdTag, err := h.dbPool.Exec(c.Context(), query, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete task"})
	}

	if cmdTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
