package database

import (
	"context"
	"fmt"
	"log"

	"github.com/NERFTHISPLS/rest-todo-list/internal/config"
	"github.com/jackc/pgx/v5"
)

const fmtDSN = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

func Setup(cfg *config.ConfDB) error {
	dsn := fmt.Sprintf(fmtDSN, cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Printf("unable to connect to the database: %s\n", err)
		return err
	}
	defer conn.Close(context.Background())

	if err := conn.Ping(context.Background()); err != nil {
		log.Printf("error making a test ping to the server: %s\n", err)
		return err
	}

	return nil
}
