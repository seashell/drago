package boltdb

import (
	"errors"

	"github.com/seashell/drago/drago/domain"
)

type NetworkRepositoryAdapter struct {
	backend *Backend
}

func NewNetworkRepositoryAdapter(backend *Backend) domain.NetworkRepository {
	return &NetworkRepositoryAdapter{backend}
}

func (a *NetworkRepositoryAdapter) GetByID(id string) (*domain.Network, error) {
	return nil, errors.New("not implemented")
}

func (a *NetworkRepositoryAdapter) Create(n *domain.Network) (*string, error) {
	return nil, errors.New("not implemented")
}

func (a *NetworkRepositoryAdapter) Update(n *domain.Network) (*string, error) {
	return nil, errors.New("not implemented")
}

func (a *NetworkRepositoryAdapter) DeleteByID(id string) (*string, error) {
	return nil, errors.New("not implemented")
}

func (a *NetworkRepositoryAdapter) FindAll(pageInfo domain.PageInfo) ([]*domain.Network, *domain.Page, error) {
	return nil, nil, errors.New("not implemented")
}
