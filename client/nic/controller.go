package nic

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os/exec"
	"time"

	structs "github.com/seashell/drago/drago/structs"
	netlink "github.com/vishvananda/netlink"
	wgctrl "golang.zx2c4.com/wireguard/wgctrl"
	wgtypes "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const (
	linkTypeWireguard = "wireguard"
)

// Config contains configurations for a network controller.
type Config struct {
	// InterfacePrefix defines the string prepended to each interface name.
	// Example: if InterfacePrefix is "abc", interfaces will be named as
	// "abc-xxxxxx", where "xxxxxx" is a random string.
	InterfacePrefix string

	// WireguardPath is the path to a userspace Wireguard binary. In case
	// it is not defined, Drago will try to use the kernel module.
	WireguardPath string
}

// Controller : network interface controller.
type Controller struct {
	config *Config
	key    *wgtypes.Key
	wg     *wgctrl.Client
}

// NewController :
func NewController(config *Config) (*Controller, error) {

	wg, err := wgctrl.New()
	if err != nil {
		return nil, err
	}

	pk, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	c := &Controller{
		config: config,
		key:    &pk,
		wg:     wg,
	}

	return c, nil
}

// ListInterfaces returns a slice of all network interfaces managed by
// the controller.
func (c *Controller) ListInterfaces() ([]*structs.Interface, error) {

	out := []*structs.Interface{}

	links, err := listLinksWithPrefix(c.config.InterfacePrefix)
	if err != nil {
		return nil, err
	}

	for _, l := range links {
		out = append(out, &structs.Interface{
			Name: l.Attrs().Name,
			// TODO: populate other fields
		})
	}

	return out, nil
}

// DeleteInterface deletes a network interface and all associated routes.
func (c *Controller) DeleteInterface(name string) error {
	err := deleteLinkAndRoutes(name)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAllInterfaces deletes all network interfaces and routes.
func (c *Controller) DeleteAllInterfaces() error {
	err := deleteLinkAndRoutesWithPrefix(c.config.InterfacePrefix)
	if err != nil {
		return err
	}
	return nil
}

// PrivateKey returns ...
func (c *Controller) PrivateKey() *wgtypes.Key {
	return c.key
}

// CreateInterface ...
func (c *Controller) CreateInterface(iface *structs.Interface) error {

	linkName, linkAlias := c.randomInterfaceName(), iface.Name

	err := c.createLink(linkName, linkAlias)
	if err != nil {
		return err
	}

	// Apply WireGuard configurations to the newly created link
	config := wgtypes.Config{
		PrivateKey:   c.key,
		ListenPort:   nil,
		Peers:        []wgtypes.PeerConfig{},
		ReplacePeers: true,
	}

	if iface.ListenPort != 0 {
		config.ListenPort = &iface.ListenPort
	}

	for _, peer := range iface.Peers {
		p, err := c.newPeerConfig(peer)
		if err != nil {
			return err
		}
		config.Peers = append(config.Peers, *p)
	}

	err = c.wg.ConfigureDevice(linkName, config)
	if err != nil {
		return err
	}

	// Assign IP address in CIDR format to the newly created link
	err = setLinkAddress(linkName, iface.Address)
	if err != nil {
		return err
	}

	err = enableLink(linkName)
	if err != nil {
		return err
	}

	idx, err := getLinkIndex(linkName)
	if err != nil {
		return err
	}

	for _, peer := range config.Peers {
		for _, ip := range peer.AllowedIPs {
			if err = netlink.RouteAdd(&netlink.Route{LinkIndex: idx, Dst: &ip}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Controller) createLink(name string, alias string) error {

	attrs := netlink.NewLinkAttrs()
	attrs.Name, attrs.Alias = name, alias

	// Create a new WireGuard interface. If a path to a valid userspace WireGuard binary
	// was specified in the configurations, use if. Otherwise, create the interface manually.
	if c.config.WireguardPath != "" {
		err := exec.Command(c.config.WireguardPath, name).Run()
		if err != nil {
			return err
		}
	} else {
		if err := netlink.LinkAdd(&netlink.Wireguard{LinkAttrs: attrs}); err != nil {
			return err
		}
	}

	err := setLinkAlias(name, alias)
	if err != nil {
		return err
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

	if peer.PublicKey != "" {
		var key wgtypes.Key
		if key, err = wgtypes.ParseKey(peer.PublicKey); err != nil {
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

	if peer.PersistentKeepalive != 0 {
		pk := time.Duration(peer.PersistentKeepalive) * time.Second
		config.PersistentKeepaliveInterval = &pk
	}

	if peer.Address != "" {
		config.Endpoint = &net.UDPAddr{
			IP:   net.ParseIP(peer.Address),
			Port: peer.Port,
		}
	}

	return config, nil
}

func (c *Controller) randomInterfaceName() string {
	buf := make([]byte, 3)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}
	return c.config.InterfacePrefix + "-" + hex.EncodeToString(buf)
}
