package inmem

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/adapter/repository"
)

type InMemBackend struct {
	db *sqlx.DB
}

func NewInmemBackend() (*InMemBackend, error) {
	db, err := sqlx.Connect("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return nil, errors.Wrap(err, "Error creating new in-memory backend")
	}

	return &InMemBackend{db}, nil
}

func (b *InMemBackend) DB() interface{} {
	return b.db
}

func (b *InMemBackend) Type() repository.BackendType {
	return repository.BackendInMemory
}

func (b *InMemBackend) ApplyMigrations(migrations ...string) error {
	for _, m := range migrations {
		b.db.MustExec(m)
	}
	return nil
}
