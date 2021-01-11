package state

import (
	"log"
	"sync"

	"github.com/seashell/drago/drago/structs"
)

// MemDB implements a StateDB that stores data in memory and should only be
// used for testing. All methods are safe for concurrent use.
type MemDB struct {
	logger log.Logger

	// interface_id -> value
	interfaces map[string]*structs.Interface

	// peers_id -> value
	peers map[string]*structs.Peer

	mu sync.RWMutex
}

func NewMemDB(logger log.Logger) *MemDB {
	return &MemDB{
		interfaces: make(map[string]*structs.Interface),
		logger:     logger,
	}
}

func (m *MemDB) Name() string {
	return "memdb"
}

func (m *MemDB) Interfaces() ([]*structs.Interface, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ifaces := make([]*structs.Interface, 0, len(m.interfaces))
	for _, v := range m.interfaces {
		ifaces = append(ifaces, v)
	}

	return ifaces, nil
}

func (m *MemDB) UpsertInterface(iface *structs.Interface) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.interfaces[iface.ID] = iface
	return nil
}

func (m *MemDB) Peers() ([]*structs.Peer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	peers := make([]*structs.Peer, 0, len(m.peers))
	for _, v := range m.peers {
		peers = append(peers, v)
	}

	return peers, nil
}

func (m *MemDB) UpsertPeer(peer *structs.Peer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.peers[peer.ID] = peer
	return nil
}
