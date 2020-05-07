package application

import (
	"github.com/seashell/drago/server/domain"
)

// NetworkService :
type NetworkService interface {
	GetByID(id string) (*domain.Network, error)
	Create(n *domain.Network) (*string, error)
	Update(n *domain.Network) (*string, error)
	DeleteByID(id string) error
	FindAll(pageInfo domain.PageInfo) ([]*domain.Network, *domain.Page, error)
}

type networkService struct {
	nr domain.NetworkRepository
}

// NewNetworkService :
func NewNetworkService(nr domain.NetworkRepository) (NetworkService, error) {
	return &networkService{nr}, nil
}

// GetByID :
func (ns *networkService) GetByID(id string) (*domain.Network, error) {
	return ns.nr.GetByID(id)
}

// Create :
func (ns *networkService) Create(n *domain.Network) (*string, error) {
	return ns.nr.Create(n)
}

// Update :
func (ns *networkService) Update(n *domain.Network) (*string, error) {
	network, err := ns.nr.GetByID(*n.ID)
	if err != nil {
		return nil, err
	}

	mergeNetworkUpdate(network, n)

	return ns.nr.Update(network)
}

// Delete :
func (ns *networkService) DeleteByID(id string) error {
	return ns.nr.DeleteByID(id)
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
