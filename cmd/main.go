// REST API для управления списком задач
// @title REST API Todo List
// @version 1.0
// @description REST API сервис для управления задачами

// @host localhost:8080
// @basePath /
package main

import (
	"log/slog"
	"os"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
	"github.com/NERFTHISPLS/rest-todo-list/internal/database"
	"github.com/NERFTHISPLS/rest-todo-list/internal/logger"
	"github.com/NERFTHISPLS/rest-todo-list/internal/repository"
	"github.com/NERFTHISPLS/rest-todo-list/internal/server"
)

func main() {
	cfg := config.New()

	logger.Setup(cfg)

	dbpool, err := database.New(&cfg.ConfDB)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	if err := database.Migrate(dbpool); err != nil {
		slog.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	repo := repository.NewTaskRepository(dbpool)

	if err := server.Setup(&cfg.Server, repo); err != nil {
		slog.Error("server setup failed", "error", err)
		os.Exit(1)
	}
}
