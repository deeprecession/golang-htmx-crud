package db

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func CreatePostgresDatabase(log *slog.Logger, psqlInfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("Failed to open a database")
		return nil, fmt.Errorf("Failed to open a database")
	}

	err = initPostgresScheme(db)
	if err != nil {
		log.Error("Failed to init postgres scheme")
		return nil, fmt.Errorf("Failed to init postgres scheme")
	}

	return db, nil
}

func initPostgresScheme(db *sql.DB) error {
	const init_query = `CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    title TEXT UNIQUE NOT NULL,
    is_done BOOLEAN NOT NULL DEFAULT FALSE
);`

	_, err := db.Exec(init_query)

	return err
}
