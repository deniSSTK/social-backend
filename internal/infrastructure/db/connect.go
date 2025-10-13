package db

import (
	"context"
	"os"
	"social-backend/internal/infrastructure/errors"
	"social-backend/internal/infrastructure/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() *pgxpool.Pool {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		logger.Get().Fatal(errors.EnvironmentVariableNotSet.Error() + "DB_URL")
	}

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		logger.Get().Fatal("Error connecting to database: " + err.Error())
	}

	return pool
}
