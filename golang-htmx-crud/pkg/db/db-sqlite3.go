package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func CreateSQLiteDatabase(sqlitePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", sqlitePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open a sqlite database: bad path=%v", sqlitePath)
	}

	err = initSQLiteScheme(db)
	if err != nil {
		return nil, fmt.Errorf("Failed to init SQLite scheme")
	}

	return db, nil
}

func initSQLiteScheme(db *sql.DB) error {
	const init_query = `CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    title TEXT UNIQUE NOT NULL,
    is_done BOOLEAN NOT NULL DEFAULT FALSE
);`

	_, err := db.Exec(init_query)

	return err
}
