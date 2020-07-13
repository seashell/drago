package application

import (
	"github.com/seashell/drago/server/domain"
)

// NetworkService :
type NetworkService interface {
	GetByID(id string) (*domain.Network, error)
	Create(n *domain.Network) (*domain.Network, error)
	Update(n *domain.Network) (*domain.Network, error)
	DeleteByID(id string) (*domain.Network, error)
	FindAll(pageInfo domain.PageInfo) ([]*domain.Network, *domain.Page, error)
}

type networkService struct {
	nr domain.NetworkRepository
	ls domain.NetworkIPAddressLeaseService
}

// NewNetworkService :
func NewNetworkService(nr domain.NetworkRepository, ls domain.NetworkIPAddressLeaseService) (NetworkService, error) {
	return &networkService{nr, ls}, nil
}

// GetByID :
func (ns *networkService) GetByID(id string) (*domain.Network, error) {
	return ns.nr.GetByID(id)
}

// Create :
func (ns *networkService) Create(n *domain.Network) (*domain.Network, error) {
	id, err := ns.nr.Create(n)
	if err != nil {
		return nil, err
	}

	n.ID = id
	err = ns.ls.PutNetwork(n)
	if err != nil {
		return nil, err
	}

	return &domain.Network{ID: id}, nil
}

// Update :
func (ns *networkService) Update(n *domain.Network) (*domain.Network, error) {
	network, err := ns.nr.GetByID(*n.ID)
	if err != nil {
		return nil, err
	}

	mergeNetworkUpdate(network, n)

	id, err := ns.nr.Update(network)
	if err != nil {
		return nil, err
	}

	return &domain.Network{ID: id}, nil
}

// Delete :
func (ns *networkService) DeleteByID(id string) (*domain.Network, error) {
	_id, err := ns.nr.DeleteByID(id)
	if err != nil {
		return nil, err
	}

	res := &domain.Network{ID: _id}
	ns.ls.PopNetwork(res)

	return res, nil
}

// FindAllByNetworkID :
func (ns *networkService) FindAll(pageInfo domain.PageInfo) ([]*domain.Network, *domain.Page, error) {
	return ns.nr.FindAll(pageInfo)
}

func mergeNetworkUpdate(current, update *domain.Network) {
	if update.Name != nil {
		current.Name = update.Name
	}

	if update.IPAddressRange != nil {
		current.IPAddressRange = update.IPAddressRange
	}
}
