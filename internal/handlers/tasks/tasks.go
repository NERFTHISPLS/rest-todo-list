package tasks

import (
	"errors"
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

func NewHandler(repo *repository.TaskRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) List(c *fiber.Ctx) error {
	tasks, err := h.repo.List(c)
	if err != nil {
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

	if err := h.repo.Create(c, task); err != nil {
		return jsonError(c, fiber.StatusInternalServerError, "failed to create task")
	}

	return c.JSON(task)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	updates := map[string]any{}
	if err := c.BodyParser(&updates); err != nil {
		return jsonError(c, fiber.StatusBadRequest, "invalid request")
	}

	for k := range updates {
		if !allowedToUpdate[k] {
			delete(updates, k)
		}
	}

	if title, ok := updates["title"]; ok {
		str, ok := title.(string)
		if !ok || strings.TrimSpace(str) == "" {
			return jsonError(c, fiber.StatusBadRequest, "title cannot be empty")
		}
	}

	if status, ok := updates["status"]; ok {
		str, ok := status.(string)
		if !ok || !isValidStatus(str) {
			return jsonError(c, fiber.StatusBadRequest, "invalid status")
		}
	}

	t, err := h.repo.Update(c, id, updates)
	if err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			return jsonError(c, fiber.StatusNotFound, "task not found")
		}
		return jsonError(c, fiber.StatusInternalServerError, "failed to update task")
	}

	return c.JSON(t)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	if err := h.repo.Delete(c, id); err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			return jsonError(c, fiber.StatusNotFound, "task not found")
		}
		return jsonError(c, fiber.StatusInternalServerError, "failed to delete task")
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

func isValidStatus(s string) bool {
	return s == "new" || s == "in_progress" || s == "done"
}
