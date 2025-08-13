package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(dbpool *pgxpool.Pool) error {
	slog.Info("starting database migration")

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
		slog.Error("error while creating the database", "error", err)
		return err
	}

	slog.Info("migration completed successfully", "table", "tasks")
	return nil
}
