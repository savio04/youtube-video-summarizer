package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/savio04/youtube-video-summarizer/internal/env"
)

var Db *pgxpool.Pool

func Init() error {
	username := env.GetEnvOrDie("POSTGRES_USER")
	password := env.GetEnvOrDie("POSTGRES_PASSWORD")
	database := env.GetEnvOrDie("POSTGRES_DB")
	host := env.GetEnvOrDie("POSTGRES_HOST")

	DATABASE_URL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", username, password, host, database)

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	connection, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	Db = connection

	return nil
}
