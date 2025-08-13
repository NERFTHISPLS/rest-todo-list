package tasks

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/NERFTHISPLS/rest-todo-list/internal/models"
	"github.com/NERFTHISPLS/rest-todo-list/internal/repository"
	"github.com/gofiber/fiber/v2"
)

var allowedToUpdate = map[string]bool{"title": true, "description": true, "status": true}

type Handler struct {
	repo *repository.TaskRepository
}

type taskRequest struct {
	Title       string  `json:"title" example:"Купить молоко"`
	Description *string `json:"description,omitempty" example:"Взять 2 литра и хлеб"`
	Status      string  `json:"status" example:"new"`
}

func NewHandler(repo *repository.TaskRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// List возвращает список всех задач
// @Summary Получить список всех задач
// @Description Возвращает список всех задач
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {array} models.Task "Список задач"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /tasks [get]
func (h *Handler) List(c *fiber.Ctx) error {
	ctx := c.Context()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("handling list tasks request", "ip", c.IP(), "user_agent", c.Get("User-Agent"))
	}

	tasks, err := h.repo.List(c)
	if err != nil {
		slog.Error("failed to list tasks", "error", err, "ip", c.IP())
		return jsonError(c, fiber.StatusInternalServerError, err.Error())
	}

	slog.Info("tasks listed successfully", "count", len(tasks), "ip", c.IP())

	return c.JSON(tasks)
}

// Create создает новую задачу
// @Summary Создать новую задачу
// @Description Создает новую задачу с указанными параметрами
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body taskRequest true "Данные задачи (title, description, status)"
// @Success 200 {object} models.Task "Созданная задача"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /tasks [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("handling create task request", "ip", c.IP(), "user_agent", c.Get("User-Agent"))
	}

	task := &models.Task{}
	if err := c.BodyParser(task); err != nil {
		slog.Warn("failed to parse request body", "error", err, "ip", c.IP())
		return jsonError(c, fiber.StatusBadRequest, "invalid request")
	}

	if strings.TrimSpace(task.Title) == "" {
		slog.Warn("task creation rejected: empty title", "ip", c.IP())
		return jsonError(c, fiber.StatusBadRequest, "title is required")
	}

	slog.Info("creating task", "title", task.Title, "status", task.Status, "ip", c.IP())

	if err := h.repo.Create(c, task); err != nil {
		slog.Error("failed to create task in database", "error", err, "title", task.Title, "ip", c.IP())
		return jsonError(c, fiber.StatusInternalServerError, "failed to create task")
	}

	slog.Info("task created successfully", "id", task.ID, "title", task.Title, "ip", c.IP())

	return c.JSON(task)
}

// Update обновляет существующую задачу
// @Summary Обновить задачу
// @Description Обновляет существующую задачу по ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param task body taskRequest true "Данные задачи (title, description, status)"
// @Success 200 {object} models.Task "Обновленная задача"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Задача не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /tasks/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	ctx := c.Context()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("handling update task request", "ip", c.IP(), "user_agent", c.Get("User-Agent"))
	}

	id, err := parseID(c)
	if err != nil {
		slog.Warn("invalid task ID in update request", "error", err, "ip", c.IP())
		return err
	}

	updates := map[string]any{}
	if err := c.BodyParser(&updates); err != nil {
		slog.Warn("failed to parse update request body", "error", err, "task_id", id, "ip", c.IP())
		return jsonError(c, fiber.StatusBadRequest, "invalid request")
	}

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("updating task", "id", id, "updates", updates, "ip", c.IP())
	} else {
		slog.Info("updating task", "id", id, "ip", c.IP())
	}

	for k := range updates {
		if !allowedToUpdate[k] {
			delete(updates, k)
		}
	}

	if title, ok := updates["title"]; ok {
		str, ok := title.(string)
		if !ok || strings.TrimSpace(str) == "" {
			slog.Warn("update rejected: empty title", "task_id", id, "ip", c.IP())
			return jsonError(c, fiber.StatusBadRequest, "title cannot be empty")
		}
	}

	if status, ok := updates["status"]; ok {
		str, ok := status.(string)
		if !ok || !isValidStatus(str) {
			slog.Warn("update rejected: invalid status", "task_id", id, "status", status, "ip", c.IP())
			return jsonError(c, fiber.StatusBadRequest, "invalid status")
		}
	}

	t, err := h.repo.Update(c, id, updates)
	if err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			slog.Warn("task not found for update", "task_id", id, "ip", c.IP())
			return jsonError(c, fiber.StatusNotFound, "task not found")
		}
		slog.Error("failed to update task in database", "error", err, "task_id", id, "ip", c.IP())
		return jsonError(c, fiber.StatusInternalServerError, "failed to update task")
	}

	slog.Info("task updated successfully", "id", id, "ip", c.IP())

	return c.JSON(t)
}

// Delete удаляет задачу по ID
// @Summary Удалить задачу
// @Description Удаляет задачу по указанному ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 204 "Задача успешно удалена"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Задача не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /tasks/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	ctx := c.Context()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug("handling delete task request", "ip", c.IP(), "user_agent", "User-Agent")
	}

	id, err := parseID(c)
	if err != nil {
		slog.Warn("invalid task ID in delete request", "error", err, "ip", c.IP())
		return err
	}

	slog.Info("deleting task", "id", id, "ip", c.IP())

	if err := h.repo.Delete(c, id); err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			slog.Warn("task not found for deletion", "task_id", id, "ip", c.IP())
			return jsonError(c, fiber.StatusNotFound, "task not found")
		}
		slog.Error("failed to delete task from database", "error", err, "task_id", id, "ip", c.IP())
		return jsonError(c, fiber.StatusInternalServerError, "failed to delete task")
	}

	slog.Info("task deleted successfully", "id", id, "ip", c.IP())

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

func isValidStatus(s string) bool {
	return s == "new" || s == "in_progress" || s == "done"
}
