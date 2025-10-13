package scripts

import (
	"database/sql"
	"errors"

	"os"
	"path/filepath"
	e "social-backend/internal/infrastructure/errors"
	"social-backend/internal/infrastructure/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Migrate() {
	log := logger.Get().Sugar()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal(e.EnvironmentVariableNotSet.Error() + "DB_URL")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get cwd: %v", err)
	}
	migrationPath := filepath.Join(cwd, "internal/infrastructure/db/migrations")
	migrationUrl := "file://" + filepath.ToSlash(migrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		migrationUrl,
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to run migration: %v", err)
	}

	log.Info("Migration completed successfully")
}
