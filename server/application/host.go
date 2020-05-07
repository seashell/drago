package application

import (
	"github.com/seashell/drago/server/domain"
)

// HostService :
type HostService interface {
	GetByID(id string) (*domain.Host, error)
	Create(h *domain.Host) (*string, error)
	Update(h *domain.Host) (*string, error)
	DeleteByID(id string) error
	FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error)
}

type hostService struct {
	hr domain.HostRepository
}

// NewHostService :
func NewHostService(hr domain.HostRepository) (HostService, error) {
	return &hostService{hr}, nil
}

// GetByID :
func (hs *hostService) GetByID(id string) (*domain.Host, error) {
	return hs.hr.GetByID(id)
}

// Create :
func (hs *hostService) Create(h *domain.Host) (*string, error) {
	return hs.hr.Create(h)
}

// Update :
func (hs *hostService) Update(h *domain.Host) (*string, error) {
	host, err := hs.hr.GetByID(*h.ID)
	if err != nil {
		return nil, err
	}

	mergeHostUpdate(host, h)

	return hs.hr.Update(host)
}

// Delete :
func (hs *hostService) DeleteByID(id string) error {
	return hs.hr.DeleteByID(id)
}

// FindAllByNetworkID :
func (hs *hostService) FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error) {
	return hs.hr.FindAllByNetworkID(id, pageInfo)
}

func mergeHostUpdate(current, update *domain.Host) {
	if update.Name != nil {
		current.Name = update.Name
	}

	if update.IPAddress != nil {
		current.IPAddress = update.IPAddress
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
