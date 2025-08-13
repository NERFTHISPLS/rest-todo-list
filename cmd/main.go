package main

import (
	"log"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
	"github.com/NERFTHISPLS/rest-todo-list/internal/database"
	"github.com/NERFTHISPLS/rest-todo-list/internal/server"
)

func main() {
	cfg := config.New()

	db, err := database.New(&cfg.ConfDB)
	if err != nil {
		log.Fatalf("error while connecting to database: %s\n", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("error while migrating to database: %s\n", err)
	}

	if err := server.Setup(&cfg.Server); err != nil {
		log.Fatalf("server setup failed: %s\n", err)
	}
}
