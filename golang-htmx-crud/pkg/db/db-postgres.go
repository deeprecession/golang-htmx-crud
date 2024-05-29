package db

import (
	"database/sql"
	"fmt"
	"net"
	"os"
)

func CreatePostgresDatabase(psqlInfo string) (*sql.DB, error) {
	postgresDB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open a database: %w", err)
	}

	err = initPostgresScheme(postgresDB)
	if err != nil {
		return nil, fmt.Errorf("failed to open a database: %w", err)
	}

	return postgresDB, nil
}

func initPostgresScheme(postgresDB *sql.DB) error {
	const initQuery = `
CREATE TABLE IF NOT EXISTS task (
    id SERIAL PRIMARY KEY NOT NULL,
    title TEXT UNIQUE NOT NULL,
    is_done BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY NOT NULL,
    login TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_task (
    user_id INT NOT NULL,
    task_id INT NOT NULL,

    FOREIGN KEY (user_id) REFERENCES "user"(id),
    FOREIGN KEY (task_id) REFERENCES task(id),

    PRIMARY KEY (user_id, task_id)
);`

	_, err := postgresDB.Exec(initQuery)
	if err != nil {
		return fmt.Errorf("failed to init postgres db: %w", err)
	}

	return nil
}

func GetPsqlInfoFromEnv() (string, error) {
	const funcErrMsg = "db.GetPsqlInfoFromEnv"

	username := os.Getenv("DB_USER")
	if username == "" {
		return "", fmt.Errorf("%s: DB_USER env variable is empty", funcErrMsg)
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		return "", fmt.Errorf("%s: DB_PORT env variable is empty", funcErrMsg)
	}

	pswd := os.Getenv("DB_PASSWORD")
	if pswd == "" {
		return "", fmt.Errorf("%s: DB_PASSWORD env variable is empty", funcErrMsg)
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		return "", fmt.Errorf("%s: DB_HOST env variable is empty", funcErrMsg)
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		return "", fmt.Errorf("%s: DB_NAME env variable is empty", funcErrMsg)
	}

	psqlInfo := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		username,
		pswd,
		net.JoinHostPort(host, port),
		dbname,
	)

	return psqlInfo, nil
}
