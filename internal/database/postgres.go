package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Db *pgxpool.Pool

func Init() error {
	// TODO: And env variables
	const DATABASE_URL string = "postgres://postgres:local@localhost:5432/postgres?"

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
