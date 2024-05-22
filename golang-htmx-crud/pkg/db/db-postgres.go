package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
)

func CreatePostgresDatabase(log *slog.Logger, psqlInfo string) (*sql.DB, error) {
	postgresDB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("failed to open a database")

		return nil, fmt.Errorf("failed to open a database: %w", err)
	}

	err = initPostgresScheme(postgresDB)
	if err != nil {
		log.Error("failed to init postgres scheme", "err", err)

		return nil, fmt.Errorf("failed to open a database: %w", err)
	}

	return postgresDB, nil
}

func initPostgresScheme(postgresDB *sql.DB) error {
	const initQuery = `CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY NOT NULL,
    title TEXT UNIQUE NOT NULL,
    is_done BOOLEAN NOT NULL DEFAULT FALSE
);`

	_, err := postgresDB.Exec(initQuery)
	if err != nil {
		return fmt.Errorf("failed to init postgres db: %w", err)
	}

	return nil
}

func GetPsqlInfoFromEnv() string {
	username := os.Getenv("DB_USER")
	port := os.Getenv("DB_PORT")
	pswd := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		pswd,
		host,
		port,
		dbname,
	)

	return psqlInfo
}
