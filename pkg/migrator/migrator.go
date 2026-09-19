package migrator

import (
	"database/sql"
	"log"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	DB         *sql.DB
	Migrations *migrate.FileMigrationSource
	Dialect    string
}

func NewMigrator(db *sql.DB, dialect string) *Migrator {

	migration := &migrate.FileMigrationSource{
		Dir: "/home/amir/Desktop/ShareExpense/internal/platform/migrations",
	}
	return &Migrator{
		Migrations: migration,
		Dialect:    dialect,
		DB:         db,
	}
}

func (m *Migrator) Up() {
	count, err := migrate.Exec(m.DB, m.Dialect, m.Migrations, migrate.Up)
	if err != nil {
		panic(err)
	}

	log.Printf("Applied %d migrations!", count)
}

func (m *Migrator) Down() {
	count, err := migrate.Exec(m.DB, m.Dialect, m.Migrations, migrate.Down)
	if err != nil {
		panic(err)
	}

	log.Printf("Applied %d migrations!", count)
}
