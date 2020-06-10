package postgres

import (
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/seashell/drago/server/adapter/repository"
)

type PostgreSQLBackend struct {
	db *sqlx.DB
}

func NewPostgreSQLBackend(address string, port uint16, dbname, username, password, sslmode string) (*PostgreSQLBackend, error) {

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		username,
		password,
		address,
		port,
		dbname,
		sslmode)

	db, err := sqlx.Connect("pgx", url)
	if err != nil {
		return nil, err
	}

	// Apply migrations on creation
	for _, m := range Migrations {
		_, err := db.Exec(m)
		if err != nil {
			return nil, err
		}
	}

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
