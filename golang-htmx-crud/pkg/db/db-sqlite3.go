package db

import (
	"database/sql"
	"fmt"
)

func CreateSQLiteDatabase(sqlitePath string) (*sql.DB, error) {
	sqliteDB, err := sql.Open("sqlite3", sqlitePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open a sqlite database: %w", err)
	}

	err = initSQLiteScheme(sqliteDB)
	if err != nil {
		return nil, fmt.Errorf("failed to create a sqlite db: %w", err)
	}

	return sqliteDB, nil
}

func initSQLiteScheme(litesqlDB *sql.DB) error {
	const initQuery = `CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    title TEXT UNIQUE NOT NULL,
    is_done BOOLEAN NOT NULL DEFAULT FALSE
);`

	_, err := litesqlDB.Exec(initQuery)
	if err != nil {
		return fmt.Errorf("failed to init a sqlite scheme: %w", err)
	}

	return nil
}
