package repository

import (
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

const bucketTypeHost = "host"

type inmemHostRepositoryAdapter struct {
}

// NewInmemHostRepositoryAdapter :
func NewInmemHostRepositoryAdapter(backend Backend) (domain.HostRepository, error) {
	return nil, errors.New("Not implemented")
}
