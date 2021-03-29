package nic

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	structs "github.com/seashell/drago/drago/structs"
	util "github.com/seashell/drago/pkg/util"
	netlink "github.com/vishvananda/netlink"
	wgctrl "golang.zx2c4.com/wireguard/wgctrl"
	wgtypes "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const (
	linkTypeWireguard = "wireguard"
)

// Config contains configurations for a network controller.
type Config struct {
	// InterfacesPrefix defines the string prepended to each interface name.
	// Example: if InterfacePrefix is "abc", interfaces will be named as
	// "abc-xxxxxx", where "xxxxxx" is a random string.
	InterfacesPrefix string

	// WireguardPath is the path to a userspace Wireguard binary. In case
	// it is not defined, Drago will try to use the kernel module.
	WireguardPath string

	// KeyStore is an implementation of the KeyStore interface, used by the
	// Controller to cache private keys for each interface.
	KeyStore PrivateKeyStore
}

// Controller : network interface controller.
type Controller struct {
	config *Config
	wg     *wgctrl.Client
}

// NewController :
func NewController(config *Config) (*Controller, error) {

	wg, err := wgctrl.New()
	if err != nil {
		return nil, err
	}

	if config.KeyStore == nil {
		panic("must provide a key store")
	}

	c := &Controller{
		config: config,
		wg:     wg,
	}

	return c, nil
}

// Interfaces returns a slice of all network interfaces managed by
// the controller.
func (c *Controller) Interfaces() ([]*structs.Interface, error) {

	out := []*structs.Interface{}

	links, err := linksByPrefix(c.config.InterfacesPrefix)
	if err != nil {
		return nil, err
	}

	for _, l := range links {

		dev, err := c.wg.Device(l.Attrs().Name)
		if err != nil {
			return nil, err
		}

		out = append(out, &structs.Interface{
			ID:         l.Attrs().Alias,
			Name:       util.StrToPtr(l.Attrs().Name),
			ListenPort: &dev.ListenPort,
			PublicKey:  util.StrToPtr(dev.PrivateKey.PublicKey().String()),
			// TODO: capture other information that might be useful e.g. for diagnosis (see l.Attrs().Statistics)
		})
	}

	return out, nil
}

// DeleteInterfaceByName deletes a network interface and all associated routes by name.
func (c *Controller) DeleteInterfaceByName(s string) error {
	err := deleteLinkAndRoutesByName(s)
	if err != nil {
		return err
	}
	return nil
}

// DeleteInterfaceByAlias deletes a network interface and all associated routes by alias.
func (c *Controller) DeleteInterfaceByAlias(s string) error {
	err := deleteLinksAndRoutesByAlias(s)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAllInterfaces deletes all network interfaces and routes.
func (c *Controller) DeleteAllInterfaces() error {
	err := deleteLinksAndRoutesByPrefix(c.config.InterfacesPrefix)
	if err != nil {
		return err
	}
	return nil
}

// CreateInterface ...
func (c *Controller) CreateInterface(iface *structs.Interface) error {

	linkName, linkAlias := c.randomInterfaceName(), iface.ID

	err := c.createLink(linkName, linkAlias)
	if err != nil {
		return err
	}

	return c.UpdateInterface(iface)

}

// UpdateInterface :
func (c *Controller) UpdateInterface(iface *structs.Interface) error {
	return c.configureLink(iface)
}

func (c *Controller) createLink(name string, alias string) error {

	attrs := netlink.NewLinkAttrs()
	attrs.Name, attrs.Alias = name, alias

	// Create a new WireGuard interface. If a path to a valid userspace WireGuard binary
	// was specified in the configurations, use if. Otherwise, create the interface manually.
	if c.config.WireguardPath != "" {

		wgt, err := wgImplementationType(c.config.WireguardPath)
		if err != nil {
			return err
		}

		var cmd *exec.Cmd

		switch wgt {
		case "wireguard-go":
			cmd = exec.Command(c.config.WireguardPath, name)
		case "boringtun":
			cmd = exec.Command(c.config.WireguardPath, name, "--disable-drop-privileges", "root")
		}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("can't create network interface with specified wireguard binary: %s", err.Error())
		}

	} else {
		fmt.Println("USING KERNELSPACE WIREGUARD")
		if err := netlink.LinkAdd(&netlink.Wireguard{LinkAttrs: attrs}); err != nil {
			return fmt.Errorf("can't create network interface : %s", err.Error())
		}
	}

	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	err = netlink.LinkSetAlias(link, alias)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) configureLink(iface *structs.Interface) error {

	link, err := netlink.LinkByAlias(iface.ID)
	if err != nil {
		return err
	}

	err = netlink.LinkSetDown(link)
	if err != nil {
		return err
	}

	// Create a new private key for this interface.
	// If a private key has already been created, use it.
	// TODO: implement an expiration/rotation strategy.
	key, err := c.config.KeyStore.KeyByID(iface.ID)
	if err != nil {
		wgKey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return fmt.Errorf("could not generate private key: %v", err)
		}
		key = &PrivateKey{
			ID:        iface.ID,
			Key:       wgKey.String(),
			CreatedAt: time.Now().Unix(),
		}
		err = c.config.KeyStore.UpsertKey(key)
		if err != nil {
			return fmt.Errorf("could not persist private key: %v", err)
		}
	}

	wgKey, err := wgtypes.ParseKey(key.Key)
	if err != nil {
		return fmt.Errorf("could not parse private key: %v", err)
	}

	config := wgtypes.Config{
		PrivateKey:   &wgKey,
		ListenPort:   iface.ListenPort,
		Peers:        []wgtypes.PeerConfig{},
		ReplacePeers: true,
	}

	for _, peer := range iface.Peers {
		peerConfig, err := c.newPeerConfig(peer)
		if err != nil {
			return err
		}
		config.Peers = append(config.Peers, *peerConfig)
	}

	err = c.wg.ConfigureDevice(link.Attrs().Name, config)
	if err != nil {
		return err
	}

	// Assign IP address in CIDR format to the newly created link
	err = setLinkAddress(link, iface.Address)
	if err != nil {
		return err
	}

	err = netlink.LinkSetUp(link)
	if err != nil {
		return err
	}

	for _, peerConfig := range config.Peers {
		for _, ip := range peerConfig.AllowedIPs {
			if err = netlink.RouteReplace(&netlink.Route{LinkIndex: link.Attrs().Index, Dst: &ip}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Controller) newPeerConfig(peer *structs.Peer) (*wgtypes.PeerConfig, error) {

	var err error

	config := &wgtypes.PeerConfig{
		Remove:                      false,
		UpdateOnly:                  false,
		ReplaceAllowedIPs:           true,
		PresharedKey:                nil,
		PublicKey:                   wgtypes.Key{},
		AllowedIPs:                  []net.IPNet{},
		Endpoint:                    nil,
		PersistentKeepaliveInterval: nil,
	}

	if peer.PublicKey != nil {
		var key wgtypes.Key
		if key, err = wgtypes.ParseKey(*peer.PublicKey); err != nil {
			return nil, err
		}
		config.PublicKey = key
	}

	for _, ip := range peer.AllowedIPs {
		_, parsed, err := net.ParseCIDR(ip)
		if err != nil {
			return nil, err
		}
		config.AllowedIPs = append(config.AllowedIPs, *parsed)
	}

	if peer.PersistentKeepalive != nil {
		pk := time.Duration(*peer.PersistentKeepalive) * time.Second
		config.PersistentKeepaliveInterval = &pk
	}

	if peer.Address != nil {
		port := 51820 // TODO: should we use this or the zero value by default?
		if peer.Port != nil {
			port = *peer.Port
		}
		config.Endpoint = &net.UDPAddr{
			IP:   net.ParseIP(*peer.Address),
			Port: port,
		}
	}

	return config, nil
}

func (c *Controller) randomInterfaceName() string {
	buf := make([]byte, 3)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}
	return c.config.InterfacesPrefix + "-" + hex.EncodeToString(buf)
}

func wgImplementationType(path string) (string, error) {

	var out bytes.Buffer

	cmd := exec.Command(path, "--version")
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	if strings.Contains(out.String(), "wireguard-go") {
		return "wireguard-go", nil
	} else if strings.Contains(out.String(), "boringtun") {
		return "boringtun", nil
	}

	return "", fmt.Errorf("unknown wireguard implementation")
}
