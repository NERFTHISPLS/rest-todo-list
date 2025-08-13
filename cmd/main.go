package main

import (
	"log"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
	"github.com/NERFTHISPLS/rest-todo-list/internal/database"
	"github.com/NERFTHISPLS/rest-todo-list/internal/repository"
	"github.com/NERFTHISPLS/rest-todo-list/internal/server"
)

func main() {
	cfg := config.New()

	dbpool, err := database.New(&cfg.ConfDB)
	if err != nil {
		log.Fatalf("error while connecting to database: %s\n", err)
	}
	defer dbpool.Close()

	if err := database.Migrate(dbpool); err != nil {
		log.Fatalf("error while migrating to database: %s\n", err)
	}

	repo := repository.NewTaskRepository(dbpool)

	if err := server.Setup(&cfg.Server, repo); err != nil {
		log.Fatalf("server setup failed: %s\n", err)
	}
}
