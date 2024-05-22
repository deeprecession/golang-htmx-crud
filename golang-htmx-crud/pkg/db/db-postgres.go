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

	err = db.Ping()
	if err != nil {
		log.Error("Failed to ping a database after opening it")
		return nil, fmt.Errorf("Failed to ping a database after opening it")
	}

	return db, nil
}
