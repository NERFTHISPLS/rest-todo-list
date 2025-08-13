package tasks

import (
	"errors"
	"fmt"
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
		return jsonError(c, fiber.StatusInternalServerError, "database connection is not initialized")
	}

	query := `SELECT * FROM tasks`
	rows, err := h.dbPool.Query(c.Context(), query)
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task

		if err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return jsonError(c, fiber.StatusInternalServerError, err.Error())
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(tasks)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	task := &models.Task{}
	if err := c.BodyParser(task); err != nil {
		return jsonError(c, fiber.StatusBadRequest, "invalid request")
	}

	if strings.TrimSpace(task.Title) == "" {
		return jsonError(c, fiber.StatusBadRequest, "title is required")
	}

	if strings.TrimSpace(task.Status) == "" {
		task.Status = models.DefaultTaskStatus
	}

	query := `INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id`
	if err := h.dbPool.QueryRow(
		c.Context(), query, task.Title, task.Description, task.Status,
	).Scan(&task.ID); err != nil {
		return jsonError(c, fiber.StatusInternalServerError, "failed to create task")
	}

	return c.Status(fiber.StatusOK).JSON(task)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	task := &models.Task{}
	if err := c.BodyParser(task); err != nil {
		return jsonError(c, fiber.StatusBadRequest, "invalid request")
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
		return jsonError(c, fiber.StatusBadRequest, "no fields to update")
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
			return jsonError(c, fiber.StatusNotFound, "task not found")
		}

		return jsonError(c, fiber.StatusInternalServerError, "failed to update task")
	}

	return c.Status(fiber.StatusOK).JSON(updatedTask)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	query := `DELETE FROM tasks WHERE id = $1`
	cmdTag, err := h.dbPool.Exec(c.Context(), query, id)
	if err != nil {
		return jsonError(c, fiber.StatusInternalServerError, "failed to delete task")
	}

	if cmdTag.RowsAffected() == 0 {
		return jsonError(c, fiber.StatusNotFound, "task not found")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func parseID(c *fiber.Ctx) (int, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	return id, nil
}

func jsonError(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{"error": msg})
}
