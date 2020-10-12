package application

import (
	"context"
	"fmt"

	domain "github.com/seashell/drago/client/domain"
	structs "github.com/seashell/drago/drago/application/structs"
)

// ReconciliationService ...
type ReconciliationService interface {
	Reconcile(ctx context.Context) error
}

type reconciliationService struct {
	config *Config
}

// NewReconciliationService ...
func NewReconciliationService(config *Config) ReconciliationService {

	s := &reconciliationService{
		config: config,
	}

	state, _ := s.config.StateRepository.Get()

	if err := s.reconcileState(state); err != nil {
		panic(err)
	}

	return s
}
func (s *reconciliationService) Reconcile(ctx context.Context) error {

	current := s.currentState()

	in := &structs.HostSynchronizeInput{
		Host: structs.Host{
			Interfaces: []*structs.Interface{},
			Peers:      []*structs.Peer{},
		},
	}

	for _, link := range current.Interfaces {
		fmt.Println(link)
	}

	for _, peer := range current.Peers {
		fmt.Println(peer)
	}

	// Use API client to fetch desired state
	out, err := s.config.DragoGateway.Agent().SynchronizeSelf(ctx, in)
	if err != nil {
		return err
	}

	desired := &domain.Host{}

	for _, link := range out.Interfaces {
		desired.Interfaces = append(desired.Interfaces, &domain.Interface{
			Name: link.Name,
			// ...
		})
	}

	for _, peer := range out.Peers {
		fmt.Println(peer)
	}

	// Apply desired state
	s.reconcileState(desired)

	return nil
}

func (s *reconciliationService) currentState() *domain.Host {
	return &domain.Host{
		Interfaces: []*domain.Interface{},
		Peers:      []*domain.Peer{},
	}
}

func (s *reconciliationService) reconcileState(desired *domain.Host) error {

	// TODO: avoid deleting all interfaces on every reconciliation
	err := s.config.InterfaceController.DeleteAllInterfaces()
	if err != nil {
		return err
	}

	for _, iface := range desired.Interfaces {
		peers := []*domain.Peer{}
		for _, peer := range desired.Peers {
			if iface.Name == peer.Interface {
				peers = append(peers, peer)
			}
		}
		iface.Peers = peers
		err := s.config.InterfaceController.CreateInterface(iface)
		if err != nil {
			return err
		}
	}

	return nil
}
