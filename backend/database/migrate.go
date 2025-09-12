package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func RunMigrations(dsn string, migrationsPath string, logger *slog.Logger) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres.WithInstance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance: %w", err)
	}

	err = m.Up()
	if err != nil {
		var dirtyErr migrate.ErrDirty
		if errors.As(err, &dirtyErr) {
			logger.Warn("Database is dirty. Forcing version", "version", dirtyErr.Version)
			if forceErr := m.Force(int(dirtyErr.Version)); forceErr != nil {
				return fmt.Errorf("force version failed: %w", forceErr)
			}
			err = m.Up()
			if err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("migrate up after force: %w", err)
			}
		} else if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migrate up: %w", err)
		}
	}

	logger.Info("Migrations applied successfully")

	return nil
}
