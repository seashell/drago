package nic

import (
	structs "github.com/seashell/drago/drago/structs"
)

// Key resolver func returns the private key associated to the ID passed
// as argument, so that we can decouple key storage from parsing and utilization.
type KeyResolverFunc func(id string) string

// NetworkInterfaceController provides network configuration capabilities.
type NetworkInterfaceController interface {
	GenerateKey() (string, error)
	CreateInterfaceWithKey(iface *structs.Interface, k string) error
	InterfacesWithPublicKey(f KeyResolverFunc) ([]*structs.Interface, error)
	Interfaces() ([]*structs.Interface, error)
	DeleteInterfaceByAlias(s string) error
	DeleteInterfaceByName(s string) error
	DeleteAllInterfaces() error
}
