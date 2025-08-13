package repository

import (
	"errors"
	"fmt"
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
	query := `
		SELECT id, title, COALESCE(description, ''), status, created_at, updated_at
		FROM tasks`

	rows, err := r.dbPool.Query(c.Context(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}

	for rows.Next() {
		var t models.Task
		var desc string

		if err := rows.Scan(&t.ID, &t.Title, &desc, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}

		if desc != "" {
			t.Description = &desc
		} else {
			t.Description = nil
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepository) Create(c *fiber.Ctx, task *models.Task) error {
	if task.Status == "" {
		task.Status = models.DefaultTaskStatus
	}

	query := `
		INSERT INTO tasks (title, description, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.dbPool.QueryRow(
		c.Context(),
		query,
		task.Title,
		task.Description,
		task.Status,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (r *TaskRepository) Update(c *fiber.Ctx, id int, updates map[string]any) (*models.Task, error) {
	if len(updates) == 0 {
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

	row := r.dbPool.QueryRow(c.Context(), query, args...)

	t := &models.Task{}
	var desc string
	if err := row.Scan(&t.ID, &t.Title, &desc, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			return nil, fiber.ErrNotFound
		}

		return nil, err
	}

	if desc != "" {
		t.Description = &desc
	} else {
		t.Description = nil
	}

	return t, nil
}

func (r *TaskRepository) Delete(c *fiber.Ctx, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	cmd, err := r.dbPool.Exec(c.Context(), query, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fiber.ErrNotFound
	}

	return nil
}
