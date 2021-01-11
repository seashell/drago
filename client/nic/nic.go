package nic

import (
	structs "github.com/seashell/drago/drago/structs"
)

// NetworkInterfaceController provides network configuration capabilities.
type NetworkInterfaceController interface {
	CreateInterface(iface *structs.Interface) error
	ListInterfaces() ([]*structs.Interface, error)
	DeleteInterface(name string) error
	DeleteAllInterfaces() error
}
