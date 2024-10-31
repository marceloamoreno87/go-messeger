package core

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DSN string
}

func (d Postgres) Connect() (db *sql.DB, err error) {
	db, err = sql.Open("postgres", d.DSN)
	if err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		return
	}

	return
}

func (d Postgres) RunMigrate() (err error) {
	migrationPath := "migrations"

	m, err := migrate.New(
		fmt.Sprintf("file:%s", migrationPath),
		d.DSN,
	)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
