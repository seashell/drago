package repository

import (
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

const bucketTypeNetwork = "network"

type inmemNetworkRepositoryAdapter struct {
}

// NewInmemNetworkRepositoryAdapter :
func NewInmemNetworkRepositoryAdapter(backend Backend) (domain.NetworkRepository, error) {
	return nil, errors.New("Not implemented")
}
