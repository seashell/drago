package application

import (
	"github.com/seashell/drago/server/domain"
)

// InterfaceService :
type InterfaceService interface {
	GetByID(id string) (*domain.Interface, error)
	Create(h *domain.Interface) (*domain.Interface, error)
	Update(h *domain.Interface) (*domain.Interface, error)
	DeleteByID(id string) (*domain.Interface, error)
	FindAll(pageInfo domain.PageInfo) ([]*domain.Interface, *domain.Page, error)
	FindAllByHostID(id string, pageInfo domain.PageInfo) ([]*domain.Interface, *domain.Page, error)
	FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Interface, *domain.Page, error)
}

type interfaceService struct {
	ifaceRepo   domain.InterfaceRepository
	networkRepo domain.NetworkRepository
}

// NewInterfaceService :
func NewInterfaceService(ir domain.InterfaceRepository, nr domain.NetworkRepository) (InterfaceService, error) {
	return &interfaceService{ir, nr}, nil
}

// GetByID :
func (s *interfaceService) GetByID(id string) (*domain.Interface, error) {
	return s.ifaceRepo.GetByID(id)
}

// Create :
func (s *interfaceService) Create(i *domain.Interface) (*domain.Interface, error) {

	// If a network ID is present, check whether the address assigned
	// to the interface lies within the allowed range
	if i.NetworkID != nil {
		if i.IPAddress != nil {
			if n, err := s.networkRepo.GetByID(*i.NetworkID); err == nil {
				err := n.CheckAddressInRange(*i.IPAddress)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	id, err := s.ifaceRepo.Create(i)
	if err != nil {
		return nil, err
	}

	return &domain.Interface{ID: id}, nil
}

// Update :
func (s *interfaceService) Update(i *domain.Interface) (*domain.Interface, error) {

	// If a network ID is present, check whether the address assigned
	// to the interface lies within the allowed range
	if i.NetworkID != nil {
		if i.IPAddress != nil {
			if n, err := s.networkRepo.GetByID(*i.NetworkID); err == nil {
				err := n.CheckAddressInRange(*i.IPAddress)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	iface, err := s.ifaceRepo.GetByID(*i.ID)
	if err != nil {
		return nil, err
	}

	mergeInterfaceUpdate(iface, i)

	id, err := s.ifaceRepo.Update(iface)
	if err != nil {
		return nil, err
	}

	return &domain.Interface{ID: id}, nil
}

// Delete :
func (s *interfaceService) DeleteByID(id string) (*domain.Interface, error) {
	_id, err := s.ifaceRepo.DeleteByID(id)
	if err != nil {
		return nil, err
	}
	return &domain.Interface{ID: _id}, nil
}

// FindAll :
func (s *interfaceService) FindAll(pageInfo domain.PageInfo) ([]*domain.Interface, *domain.Page, error) {
	return s.ifaceRepo.FindAll(pageInfo)
}

// FindAllByNetworkID :
func (s *interfaceService) FindAllByHostID(id string, pageInfo domain.PageInfo) ([]*domain.Interface, *domain.Page, error) {
	return s.ifaceRepo.FindAllByHostID(id, pageInfo)
}

// FindAllByNetworkID :
func (s *interfaceService) FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Interface, *domain.Page, error) {
	return s.ifaceRepo.FindAllByNetworkID(id, pageInfo)
}

func mergeInterfaceUpdate(current, update *domain.Interface) {
	if update.Name != nil {
		current.Name = update.Name
	}

	if update.NetworkID != nil {
		current.NetworkID = update.NetworkID
	}

	if update.IPAddress != nil {
		current.IPAddress = update.IPAddress
	}

	if update.PublicKey != nil {
		current.PublicKey = update.PublicKey
	}

	if update.ListenPort != nil {
		current.ListenPort = update.ListenPort
	}

	if update.PublicKey != nil {
		current.PublicKey = update.PublicKey
	}

	if update.Table != nil {
		current.Table = update.Table
	}

	if update.DNS != nil {
		current.DNS = update.DNS
	}

	if update.MTU != nil {
		current.MTU = update.MTU
	}

	if update.PreUp != nil {
		current.PreUp = update.PreUp
	}

	if update.PostUp != nil {
		current.PostUp = update.PostUp
	}

	if update.PreDown != nil {
		current.PreDown = update.PreDown
	}

	if update.PostDown != nil {
		current.PostDown = update.PostDown
	}
}
