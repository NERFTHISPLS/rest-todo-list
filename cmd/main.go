package main

import (
	"log"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
	"github.com/NERFTHISPLS/rest-todo-list/internal/database"
	"github.com/NERFTHISPLS/rest-todo-list/internal/server"
)

func main() {
	cfg := config.New()

	if err := server.Setup(&cfg.Server); err != nil {
		log.Fatalf("server setup failed: %s\n", err)
	}

	log.Println("server is working")

	if err := database.Setup(&cfg.ConfDB); err != nil {
		log.Fatalf("error while connecting to database: %s\n", err)
	}

	log.Println("database is connected")
}
