package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const fmtDSN = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

func New(cfg *config.ConfDB) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(fmtDSN, cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	slog.Info("connecting to database", "host", cfg.Host, "port", cfg.Port, "database", cfg.Name)

	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		slog.Error("unable to create connection pool", "error", err)
		return nil, err
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		slog.Error("error making a test ping to the database", "error", err)
		return nil, err
	}

	slog.Info("database connected successfully")
	return dbpool, nil
}
