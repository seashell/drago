package nic

import (
	structs "github.com/seashell/drago/drago/structs"
)

// NetworkInterfaceController provides network configuration capabilities.
type NetworkInterfaceController interface {
	GenerateKey() (string, error)
	CreateInterfaceWithKey(iface *structs.Interface, k string) error
	Interfaces() ([]*structs.Interface, error)
	DeleteInterfaceByAlias(s string) error
	DeleteInterfaceByName(s string) error
	DeleteAllInterfaces() error
}
