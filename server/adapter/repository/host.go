package repository

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

func NewHostRepositoryAdapter(backend Backend) (domain.HostRepository, error) {
	switch backend.Type() {
	case BackendPostgreSQL:
		return NewPostgreSQLHostRepositoryAdapter(backend)

	case BackendInMemory:
		return NewInmemHostRepositoryAdapter(backend)

	default:
		return nil, errors.New(fmt.Sprintf("Error creating adapter for backend of type %s", backend.Type()))
	}
}
