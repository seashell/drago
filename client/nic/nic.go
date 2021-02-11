package nic

import (
	structs "github.com/seashell/drago/drago/structs"
)

// NetworkInterfaceController provides network configuration capabilities.
type NetworkInterfaceController interface {
	Interfaces() ([]*structs.Interface, error)
	CreateInterface(iface *structs.Interface) error
	UpdateInterface(iface *structs.Interface) error
	DeleteInterfaceByAlias(s string) error
	DeleteInterfaceByName(s string) error
	DeleteAllInterfaces() error
}
