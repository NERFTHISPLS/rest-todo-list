package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(dbpool *pgxpool.Pool) error {
	query := `
  CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
  );`

	if _, err := dbpool.Exec(context.Background(), query); err != nil {
		log.Printf("error while creating the database: %s\n", err)
		return err
	}

	log.Println("migration completed: tasks table created")

	return nil
}
