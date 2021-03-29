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

	mu sync.RWMutex
}

func NewRepository(logger log.Logger) *Repository {
	return &Repository{
		interfaces: make(map[string]*structs.Interface),
		logger:     logger,
	}
}

func (r *Repository) Name() string {
	return "inmem"
}

func (r *Repository) Interfaces() ([]*structs.Interface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ifaces := make([]*structs.Interface, 0, len(r.interfaces))
	for _, v := range r.interfaces {
		ifaces = append(ifaces, v)
	}

	return ifaces, nil
}

func (r *Repository) UpsertInterface(iface *structs.Interface) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.interfaces[iface.ID] = iface
	return nil
}

func (r *Repository) DeleteInterfaces(ids []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, id := range ids {
		delete(r.interfaces, id)
	}

	return nil
}
