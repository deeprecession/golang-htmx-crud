package db

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func GetDB(log *slog.Logger, enviroment string) (*sql.DB, error) {
	if enviroment == "PRODUCTION" {
		psqlInfo := GetPsqlInfoFromEnv()

		postgresDB, err := CreatePostgresDatabase(log, psqlInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to create a production database: %w", err)
		}

		return postgresDB, nil
	}

	if enviroment == "DEVELOPMENT" {
		const sqlitePath = "./sqlite.db"

		sqliteDB, err := CreateSQLiteDatabase(sqlitePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create a development database: %w", err)
		}

		return sqliteDB, nil
	}

	return nil, fmt.Errorf("bad enviroment name %q", enviroment)
}
