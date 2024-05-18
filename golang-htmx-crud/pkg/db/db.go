package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

func CreateDatabase(log *slog.Logger) (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Error("Failed to parse port for database from .env file")
		return nil, fmt.Errorf("Failed to parse port for database from .env file")
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)

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

	log.Info("Connected to the Postgres Database!", "host", host, "port", port, "dbname", dbname)

	return db, nil
}

type PostgresStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) PostgresStorage {
	return PostgresStorage{db}
}

func (storage *PostgresStorage) GetTaskList() {
}
