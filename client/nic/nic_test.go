package nic

import (
	"fmt"
	"testing"

	domain "github.com/seashell/drago/client/domain"
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

	iface, err := c.newInterface(&domain.Interface{
		Name:    "t1",
		Address: "192.168.2.1/24",
		Peers:   []*domain.Peer{},
	})

	fmt.Println(iface)

	err = c.setupInterface(iface)
	if err != nil {
		t.Fatal(err)
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
