package application

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/seashell/drago/server/adapter/repository"
	"github.com/seashell/drago/server/domain"
)

// SynchronizationService :
type SynchronizationService interface {
	GetHostSettingsByID(id string) (*domain.HostSettings, error)
	UpdateHostState(id string, state *domain.HostState) (*domain.HostState, error)
	SynchronizeHost(id string, state *domain.HostState) (*domain.HostSettings, error)
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

	settings := &domain.HostSettings{
		Interfaces: []*domain.WgInterfaceSettings{},
		Peers:      []*domain.WgPeerSettings{},
	}

	pageInfo := domain.PageInfo{
		Page:    1,
		PerPage: repository.MaxQueryRows,
	}
	for {
		ifaces, page, err := s.ifaceRepo.FindAllByHostID(id, pageInfo)
		if err != nil {
			return nil, err
		}

		for _, iface := range ifaces {
			settings.Interfaces = append(settings.Interfaces, &domain.WgInterfaceSettings{
				Name:       iface.Name,
				Address:    iface.IPAddress,
				ListenPort: iface.ListenPort,
				Table:      iface.Table,
				DNS:        iface.DNS,
				MTU:        iface.MTU,
				PreUp:      iface.PreUp,
				PostUp:     iface.PostUp,
				PreDown:    iface.PreDown,
				PostDown:   iface.PostDown,
			})
		}

		if page.Page >= page.PageCount {
			break
		}
		pageInfo.Page += 1
	}

	pageInfo.Page = 1

	for {
		links, page, err := s.linkRepo.FindAllBySourceHostID(id, pageInfo)
		if err != nil {
			return nil, err
		}

		for _, link := range links {

			sourceIface, err := s.ifaceRepo.GetByID(*link.FromInterfaceID)
			if err != nil {
				return nil, err
			}

			peerIface, err := s.ifaceRepo.GetByID(*link.ToInterfaceID)
			if err != nil {
				return nil, err
			}

			peerHost, err := s.hostRepo.GetByID(*peerIface.HostID)
			if err != nil {
				return nil, err
			}

			settings.Peers = append(settings.Peers, &domain.WgPeerSettings{
				Interface:           *sourceIface.Name,
				PublicKey:           peerIface.PublicKey,
				Address:             peerHost.AdvertiseAddress,
				Port:                peerIface.ListenPort,
				AllowedIPs:          link.AllowedIPs,
				PersistentKeepalive: link.PersistentKeepalive,
			})
		}
		if page.Page >= page.PageCount {
			break
		}
		pageInfo.Page += 1
	}

	return settings, nil
}

func (s *synchronizationService) UpdateHostState(id string, state *domain.HostState) (*domain.HostState, error) {

	pageInfo := domain.PageInfo{
		Page:    1,
		PerPage: repository.MaxQueryRows,
	}

	allHostIfaces := make(map[string]*domain.Interface)

	for {
		ifaces, page, err := s.ifaceRepo.FindAllByHostID(id, pageInfo)
		if err != nil {
			return nil, err
		}
		spew.Dump(ifaces)
		for _, iface := range ifaces {
			allHostIfaces[*iface.Name] = iface
		}
		if page.Page >= page.PageCount {
			break
		}
		pageInfo.Page += 1
	}

	for _, ifaceState := range state.Interfaces {
		iface := allHostIfaces[*ifaceState.Name]
		_, err := s.ifaceRepo.Update(&domain.Interface{ID: iface.ID, Name: iface.Name, PublicKey: ifaceState.PublicKey})
		if err != nil {
			return nil, err
		}
	}

	stateOut := &domain.HostState{}

	return stateOut, nil
}

func (s *synchronizationService) SynchronizeHost(id string, state *domain.HostState) (*domain.HostSettings, error) {

	state, err := s.UpdateHostState(id, state)
	if err != nil {
		return nil, err
	}

	settings, err := s.GetHostSettingsByID(id)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
