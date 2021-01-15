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
	// InterfacesPrefix defines the string prepended to each interface name.
	// Example: if InterfacePrefix is "abc", interfaces will be named as
	// "abc-xxxxxx", where "xxxxxx" is a random string.
	InterfacesPrefix string

	// WireguardPath is the path to a userspace Wireguard binary. In case
	// it is not defined, Drago will try to use the kernel module.
	WireguardPath string
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

	c := &Controller{
		config: config,
		wg:     wg,
	}

	return c, nil
}

// GenerateKey : generates a new 32-bytes key, and return it in base64 encoding.
func (c *Controller) GenerateKey() (string, error) {
	pk, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", err
	}
	return pk.String(), nil
}

func (c *Controller) parseKey(s string) (wgtypes.Key, error) {
	k, err := wgtypes.ParseKey(s)
	if err != nil {
		return wgtypes.Key{}, err
	}
	return k, nil
}

// InterfacesWithPublicKey returns a slice of all network interfaces managed by
// the controller, together with their public key.
func (c *Controller) InterfacesWithPublicKey(keyByID KeyResolverFunc) ([]*structs.Interface, error) {

	ifaces, err := c.Interfaces()
	if err != nil {
		return nil, err
	}

	// For each Wireguard interface in the system, retrieve the private key
	// used to configure it, set the public key computed from it on the struct.
	for _, i := range ifaces {
		pubkey := keyByID(i.ID)
		k, err := c.parseKey(pubkey)
		if err != nil {
			return nil, err
		}
		i.PublicKey = k.PublicKey().String()
	}

	return ifaces, nil
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
		out = append(out, &structs.Interface{
			ID:   l.Attrs().Alias,
			Name: l.Attrs().Name,
			// TODO: capture other information that might be useful e.g. for diagnosis
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

// CreateInterfaceWithKey ...
func (c *Controller) CreateInterfaceWithKey(iface *structs.Interface, s string) error {

	key, err := c.parseKey(s)
	if err != nil {
		return fmt.Errorf("invalid key: %s", err.Error())
	}

	linkName, linkAlias := c.randomInterfaceName(), iface.ID

	err = c.createLink(linkName, linkAlias)
	if err != nil {
		return err
	}

	// Apply WireGuard configurations to the newly created link
	config := wgtypes.Config{
		PrivateKey:   &key,
		ListenPort:   nil,
		Peers:        []wgtypes.PeerConfig{},
		ReplacePeers: true,
	}

	if iface.ListenPort != 0 {
		p := int(iface.ListenPort)
		config.ListenPort = &p
	}

	for _, peer := range iface.Peers {
		p, err := c.newPeerConfig(peer)
		if err != nil {
			return err
		}
		config.Peers = append(config.Peers, *p)
	}

	link, err := netlink.LinkByAlias(linkAlias)
	if err != nil {
		return err
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

	idx := link.Attrs().Index
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
			return fmt.Errorf("can't create network interface with specified wireguard binary: %s", err.Error())
		}
	} else {
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
	return c.config.InterfacesPrefix + "-" + hex.EncodeToString(buf)
}
