package migrations

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

type DBMigrator struct {
	db            *sqlx.DB
	migrationsDir string
}

func NewDbMigrator(db *sqlx.DB, migrationDir string) *DBMigrator {
	return &DBMigrator{db: db, migrationsDir: migrationDir}
}

func (m *DBMigrator) Migrate() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(m.db.DB, m.migrationsDir); err != nil {
		return err
	}

	return nil
}
