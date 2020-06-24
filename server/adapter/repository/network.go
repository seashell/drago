package repository

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

func NewNetworkRepositoryAdapter(backend Backend) (domain.NetworkRepository, error) {
	switch backend.Type() {
	case BackendPostgreSQL:
		return NewPostgreSQLNetworkRepositoryAdapter(backend)

	case BackendInMemory:
		return NewInmemNetworkRepositoryAdapter(backend)

	default:
		return nil, errors.New(fmt.Sprintf("Error creating adapter for backend of type %s", backend.Type()))
	}
}
