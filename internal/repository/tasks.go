package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/NERFTHISPLS/rest-todo-list/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	dbPool *pgxpool.Pool
}

func NewTaskRepository(dbPool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{dbPool: dbPool}
}

func (r *TaskRepository) List(c *fiber.Ctx) ([]models.Task, error) {
	ctx := c.Context()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("executing database query: list tasks")
	}

	query := `
		SELECT id, title, COALESCE(description, ''), status, created_at, updated_at
		FROM tasks`

	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		slog.Error("database query failed: list tasks", "error", err)
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}

	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			slog.Error("failed to scan task row", "error", err)

			return nil, err
		}

		tasks = append(tasks, t)
	}

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("database query completed: list tasks", "count", len(tasks))
	}

	return tasks, nil
}

func (r *TaskRepository) Create(c *fiber.Ctx, task *models.Task) error {
	ctx := c.Context()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("executing database query: create task", "title", task.Title, "status", task.Status)
	}

	if task.Status == "" {
		task.Status = models.DefaultTaskStatus
	}

	query := `
		INSERT INTO tasks (title, description, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.dbPool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		slog.Error("database query failed: create task", "error", err, "title", task.Title)
		return err
	}

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("database query completed: create task", "id", task.ID)
	}

	return nil
}

func (r *TaskRepository) Update(c *fiber.Ctx, id int, updates map[string]any) (*models.Task, error) {
	ctx := c.Context()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("executing database query: update task", "id", id, "updates", updates)
	}

	if len(updates) == 0 {
		slog.Warn("no fields to update", "task_id", id)
		return nil, fmt.Errorf("no fields to update")
	}

	setClauses := []string{}
	args := []any{}
	i := 1
	for k, v := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, v)
		i++
	}
	setClauses = append(setClauses, "updated_at = now()")

	query := fmt.Sprintf(`
		UPDATE tasks
		SET %s
		WHERE id = $%d
		RETURNING id, title, description, status, created_at, updated_at
	`, strings.Join(setClauses, ", "), i)

	args = append(args, id)

	row := r.dbPool.QueryRow(ctx, query, args...)

	t := &models.Task{}
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			slog.Warn("task not found for update", "task_id", id)
			return nil, fiber.ErrNotFound
		}

		slog.Error("database query failed: update task", "error", err, "task_id", id)

		return nil, err
	}

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("database query completed: update task", "id", id)
	}

	return t, nil
}

func (r *TaskRepository) Delete(c *fiber.Ctx, id int) error {
	ctx := c.Context()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("executing database query: delete task", "id", id)
	}

	query := `DELETE FROM tasks WHERE id = $1`
	cmd, err := r.dbPool.Exec(ctx, query, id)
	if err != nil {
		slog.Error("database query failed: delete task", "error", err, "task_id", id)
		return err
	}

	if cmd.RowsAffected() == 0 {
		slog.Warn("no rows affected when deleting task", "task_id", id)
		return fiber.ErrNotFound
	}

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("database query completed: delete task", "id", id, "rows_affected", cmd.RowsAffected())
	}

	return nil
}
