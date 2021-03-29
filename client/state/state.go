package state

import (
	"github.com/seashell/drago/client/nic"
	"github.com/seashell/drago/drago/structs"
)

// Repository :
type Repository interface {
	Name() string

	// Client state
	Interfaces() ([]*structs.Interface, error)
	UpsertInterface(*structs.Interface) error
	DeleteInterfaces(id []string) error

	// Key store
	KeyByID(id string) (*nic.PrivateKey, error)
	UpsertKey(key *nic.PrivateKey) error
	DeleteKey(id string) error
}
