package state

import (
	"log"
	"sync"

	"github.com/seashell/drago/drago/structs"
)

type Repository struct {
	logger log.Logger

	// interface_id -> value
	interfaces map[string]*structs.Interface

	// peers_id -> value
	peers map[string]*structs.Peer

	mu sync.RWMutex
}

func NewRepository(logger log.Logger) *Repository {
	return &Repository{
		interfaces: make(map[string]*structs.Interface),
		logger:     logger,
	}
}

func (m *Repository) Name() string {
	return "inmem"
}

func (m *Repository) Interfaces() ([]*structs.Interface, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ifaces := make([]*structs.Interface, 0, len(m.interfaces))
	for _, v := range m.interfaces {
		ifaces = append(ifaces, v)
	}

	return ifaces, nil
}

func (r *Repository) InterfaceByID(iface *structs.Interface) (*structs.Interface, error) {
	return nil, nil
}

func (m *Repository) UpsertInterface(iface *structs.Interface) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.interfaces[iface.ID] = iface
	return nil
}

func (m *Repository) Peers() ([]*structs.Peer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	peers := make([]*structs.Peer, 0, len(m.peers))
	for _, v := range m.peers {
		peers = append(peers, v)
	}

	return peers, nil
}

func (m *Repository) UpsertPeer(peer *structs.Peer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.peers[peer.ID] = peer
	return nil
}
