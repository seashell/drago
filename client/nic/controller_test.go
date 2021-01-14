package nic

import (
	"fmt"
	"testing"

	structs "github.com/seashell/drago/drago/structs"
)

func TestCreateInterface(t *testing.T) {
	config := &Config{
		InterfacesPrefix: "drago",
		WireguardPath:    "./wireguard",
	}

	c, err := NewController(config)
	if err != nil {
		t.Fatal(err)
	}

	key, err := c.GenerateKey()
	if err != nil {
		t.Error()
	}

	err = c.CreateInterfaceWithKey(&structs.Interface{
		Name:    "1234567890abcd",
		Address: "192.168.2.1/24",
		Peers:   []*structs.Peer{},
	}, key)
	if err != nil {
		t.Error(err)
	}
}

func TestListInterfaces(t *testing.T) {

	config := &Config{
		InterfacesPrefix: "drago",
		WireguardPath:    "./wireguard",
	}

	c, err := NewController(config)
	if err != nil {
		t.Fatal(err)
	}

	ifaces, err := c.Interfaces()
	if err != nil {
		t.Fatal(err)
	}

	for _, i := range ifaces {
		fmt.Printf("%s (%s)", i.Name, i.Address)
	}
}

func TestDeleteAllInterfaces(t *testing.T) {
	return
	config := &Config{
		InterfacesPrefix: "drago",
		WireguardPath:    "./wireguard",
	}

	c, err := NewController(config)
	if err != nil {
		t.Fatal(err)
	}

	err = c.DeleteAllInterfaces()
	if err != nil {
		t.Fatal(err)
	}
}
