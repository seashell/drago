package application

import (
	"errors"

	"github.com/seashell/drago/server/domain"
)

// SynchronizationService :
type SynchronizationService interface {
	GetHostSettingsByID(id string) (*domain.HostSettings, error)
}

type synchronizationService struct {
	hostRepo  domain.HostRepository
	ifaceRepo domain.InterfaceRepository
	linkRepo  domain.LinkRepository
}

// NewSynchronizationService :
func NewSynchronizationService(hostRepo domain.HostRepository, ifaceRepo domain.InterfaceRepository, linkRepo domain.LinkRepository) (SynchronizationService, error) {
	return &synchronizationService{hostRepo, ifaceRepo, linkRepo}, nil
}

func (s *synchronizationService) GetHostSettingsByID(id string) (*domain.HostSettings, error) {

	//host, err := s.hostRepo.GetByID(id)
	//ifaces, err := s.ifaceRepo.FindAllByHostID(id)
	//links, err := s.linkRepo.FindAllBySourceHostID(id)

	//settings := &domain.HostSettings{}

	return nil, errors.New("not implemented")
}
