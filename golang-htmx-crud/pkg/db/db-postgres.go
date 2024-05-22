package db

import (
	"database/sql"
	"fmt"
	"log/slog"
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
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    title TEXT UNIQUE NOT NULL,
    is_done BOOLEAN NOT NULL DEFAULT FALSE
);`

	_, err := postgresDB.Exec(initQuery)
	if err != nil {
		return fmt.Errorf("failed to init postgres db: %w", err)
	}

	return nil
}
