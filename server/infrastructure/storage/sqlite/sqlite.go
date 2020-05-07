package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/adapter/repository"
)

type SQLiteBackend struct {
	db *sqlx.DB
}

func NewSQLiteBackend(filename string) (*SQLiteBackend, error) {
	db, err := sqlx.Connect("sqlite3", filename)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating new SQLite3 backend")
	}

	return &SQLiteBackend{db}, nil
}

func (b *SQLiteBackend) DB() interface{} {
	return b.db
}

func (b *SQLiteBackend) Type() repository.BackendType {
	return repository.BackendSQLite
}

func (b *SQLiteBackend) ApplyMigrations(migrations ...string) error {
	for _, m := range migrations {
		b.db.MustExec(m)
	}
	return nil
}
