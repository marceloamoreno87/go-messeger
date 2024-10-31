package core

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

/*
Definição de variáveis de erro específicas para operações com o banco de dados Postgres.
Essas variáveis são usadas para fornecer mensagens de erro detalhadas.
*/
var (
	ErrDBConnectionFailed = errors.New("postgres.connection_failed: failed to connect to the database")
	ErrDBPingFailed       = errors.New("postgres.ping_failed: failed to ping the database")
	ErrMigrationFailed    = errors.New("postgres.migration_failed: failed to run database migrations")
)

/*
Estrutura Postgres que contém a string de conexão DSN.
DSN (Data Source Name) é usado para conectar ao banco de dados Postgres.
*/
type Postgres struct {
	DSN string
}

/*
Método Connect estabelece uma conexão com o banco de dados Postgres.
Retorna um ponteiro para a conexão do banco de dados e um erro, se houver.
*/
func (d Postgres) Connect() (db *sql.DB, err error) {
	db, err = sql.Open("postgres", d.DSN)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBConnectionFailed, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBPingFailed, err)
	}

	return db, nil
}

/*
Método RunMigrate executa as migrações do banco de dados.
As migrações garantem que o esquema do banco de dados esteja atualizado.
Retorna um erro, se houver.
*/
func (d Postgres) RunMigrate() (err error) {
	migrationPath := "migrations"

	m, err := migrate.New(
		fmt.Sprintf("file:%s", migrationPath),
		d.DSN,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMigrationFailed, err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%w: %v", ErrMigrationFailed, err)
	}

	return nil
}
