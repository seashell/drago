package nic

import (
	"fmt"
	"testing"

	structs "github.com/seashell/drago/drago/structs"
)

func TestCreateInterface(t *testing.T) {
	config := &Config{
		InterfacePrefix: "drago",
		WireguardPath:   "./wireguard",
	}

	c, err := NewController(config)
	if err != nil {
		t.Fatal(err)
	}

	err = c.CreateInterface(&structs.Interface{
		Name:    "t1",
		Address: "192.168.2.1/24",
		Peers:   []*structs.Peer{},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteAllInterfaces(t *testing.T) {
	config := &Config{
		InterfacePrefix: "drago",
		WireguardPath:   "./wireguard",
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

func TestListInterfaces(t *testing.T) {

	config := &Config{
		InterfacePrefix: "drago",
		WireguardPath:   "./wireguard",
	}

	c, err := NewController(config)
	if err != nil {
		t.Fatal(err)
	}

	ifaces, err := c.ListInterfaces()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(ifaces)
}
