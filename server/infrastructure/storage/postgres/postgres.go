package postgres

import (
	"errors"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/seashell/drago/server/adapter/repository"
)

type PostgreSQLBackend struct {
	db *sqlx.DB
}

func NewPostgreSQLBackend(address string, port uint16, dbname, username, password, sslmode string) (*PostgreSQLBackend, error) {

	pgconf := pgx.ConnConfig{
		Host:      address,
		Port:      port,
		Database:  dbname,
		User:      username,
		Password:  password,
		TLSConfig: nil,
	}

	db := sqlx.NewDb(stdlib.OpenDB(pgconf), "pgx")
	if db.Ping() != nil {
		return nil, errors.New("Error creating new PostgreSQL backend")
	}

	return &PostgreSQLBackend{db}, nil
}

func (b *PostgreSQLBackend) DB() interface{} {
	return b.db
}

func (b *PostgreSQLBackend) Type() repository.BackendType {
	return repository.BackendPostgreSQL
}

func (b *PostgreSQLBackend) ApplyMigrations(migrations ...string) error {
	for _, m := range migrations {
		b.db.MustExec(m)
	}
	return nil
}
