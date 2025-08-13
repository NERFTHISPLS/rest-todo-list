package database

import (
	"context"
	"fmt"
	"log"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const fmtDSN = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

func New(cfg *config.ConfDB) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(fmtDSN, cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Printf("unable to create connection pool: %s\n", err)
		return nil, err
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Printf("error making a test ping to the database: %s\n", err)
		return nil, err
	}

	log.Println("database is connected")

	return dbpool, err
}
