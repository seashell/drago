package repository

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

// NewInterfaceRepositoryAdapter:
func NewInterfaceRepositoryAdapter(backend Backend) (domain.InterfaceRepository, error) {
	switch backend.Type() {
	case BackendPostgreSQL:
		return NewPostgreSQLInterfaceRepositoryAdapter(backend)

	case BackendInMemory:
		return NewInmemInterfaceRepositoryAdapter(backend)

	default:
		return nil, errors.New(fmt.Sprintf("Error creating adapter for backend of type %s", backend.Type()))
	}
}
