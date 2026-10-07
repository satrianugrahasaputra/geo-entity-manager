package repository

import (
	"errors"
	"fmt"
	"log/slog"

	"geo-entity-manager/backend/migrations"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Migrate runs the embedded migrations against the given database URL.
func Migrate(databaseURL string) error {
	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("load migrations source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		return fmt.Errorf("init migration: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	slog.Info("running database migrations...")
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration up: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("get migration version: %w", err)
	}

	slog.Info("migrations completed", slog.Uint64("version", uint64(version)), slog.Bool("dirty", dirty))
	return nil
}
