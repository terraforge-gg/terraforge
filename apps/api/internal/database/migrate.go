package database

import (
	"database/sql"
	"embed"
	"log/slog"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(logger *slog.Logger, db *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {

		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		logger.Error("An error occurred running migrations", "error", err)
		return err
	}

	return nil
}
