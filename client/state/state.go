package state

import (
	"github.com/seashell/drago/drago/structs"
)

// Repository :
type Repository interface {
	Name() string
	Interfaces() ([]*structs.Interface, error)
	UpsertInterface(*structs.Interface) error
	DeleteInterfaces(id []string) error
}
