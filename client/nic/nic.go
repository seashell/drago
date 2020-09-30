package nic

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os/exec"
	"regexp"
	"time"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/seashell/drago/pkg/logger"
)

// Settings :
type Settings struct {
	Name      string
	Alias     *string
	Address   *netlink.Addr
	Wireguard *wgtypes.Config
}

// NetworkInterface :
type NetworkInterface struct {
	Settings *Settings
	Link     *netlink.Link
}

// NetworkInterfaceCtrl :
type NetworkInterfaceCtrl struct {
	NetworkInterfaces map[string]*NetworkInterface

	namePrefix   string
	wgController *wgctrl.Client
	wgPrivateKey *wgtypes.Key

	log           logger.Logger
	WireguardPath string
}

// NewCtrl :
func NewCtrl(namePrefix string, wireguardPath string, log logger.Logger) (*NetworkInterfaceCtrl, error) {

	wg, err := wgctrl.New()
	if err != nil {
		return nil, err
	}

	pk, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	return &NetworkInterfaceCtrl{
		NetworkInterfaces: make(map[string]*NetworkInterface),
		wgController:      wg,
		wgPrivateKey:      &pk,
		namePrefix:        namePrefix,
		log:               log,
		WireguardPath:     wireguardPath,
	}, nil
}

// Update :
func (n *NetworkInterfaceCtrl) Update(ts []Settings) error {
	if err := n.resetWgNetworkInterfaces(); err != nil {
		return err
	}

	for _, s := range ts {
		b := make([]byte, 5) //equals 10 charachters
		rand.Read(b)
		r := hex.EncodeToString(b)
		s.Name = n.namePrefix + r
		if err := n.ConfigureNetworkInterface(&s); err != nil {
			return err
		}
	}
	return nil
}

func (n *NetworkInterfaceCtrl) resetWgNetworkInterfaces() error {
	niList, _ := netlink.LinkList()
	for _, ni := range niList {
		//match device alias with prefix provided by n.namePrefix
		matched, err := regexp.MatchString(n.namePrefix+`.*`, ni.Attrs().Name)
		if err != nil {
			return fmt.Errorf("failed to match interface name: %v", err)
		}

		if matched {
			if err := n.DeleteNetworkInterface(ni.Attrs().Name); err != nil {
				return fmt.Errorf("failed to delete network interface: %v", err)
			}
			delete(n.NetworkInterfaces, ni.Attrs().Alias)
		}

	}
	return nil
}

// ConfigureNetworkInterface :
func (n *NetworkInterfaceCtrl) ConfigureNetworkInterface(ts *Settings) error {
	// register new interface
	lattr := netlink.NewLinkAttrs()
	lattr.Name = ts.Name
	lattr.Alias = *ts.Alias

	if n.WireguardPath != "" {
		err := exec.Command(n.WireguardPath, ts.Name).Run()
		if err != nil {
			return err
		}
	} else {
		if err := netlink.LinkAdd(&netlink.Wireguard{LinkAttrs: lattr}); err != nil {
			return fmt.Errorf("failed to create new network device: %v", err)
		}
	}

	l, err := netlink.LinkByName(ts.Name)
	if err != nil {
		return fmt.Errorf("failed to get network device by name: %v", err)
	}

	if err := netlink.LinkSetAlias(l, lattr.Alias); err != nil {
		n.log.Warnf("Setting link alias error at %v: %v\n", time.Now().Round(0), err)
	}

	// apply wireguard config
	if err := n.wgController.ConfigureDevice(ts.Name, *ts.Wireguard); err != nil {
		if err != nil {
			n.log.Warnf("Unknown wireguard configuration error at %v: %v\n", time.Now().Round(0), err)
		}
	}

	// apply new settings
	if err := netlink.AddrAdd(l, ts.Address); err != nil {
		n.log.Warnf("Adding interface address error at %v: %v\n", time.Now().Round(0), err)
	}

	if err := netlink.LinkSetUp(l); err != nil {
		n.log.Warnf("Setting interface up error at %v: %v\n", time.Now().Round(0), err)
	}

	for _, peer := range ts.Wireguard.Peers {
		for _, allowedIP := range peer.AllowedIPs {
			if err = netlink.RouteAdd(&netlink.Route{
				LinkIndex: l.Attrs().Index,
				Dst:       &allowedIP,
			}); err != nil {
				n.log.Warnf("Setting IP route error at %v: %v\n", time.Now().Round(0), err)
			}
		}
	}

	n.NetworkInterfaces[*ts.Alias] = &NetworkInterface{
		Settings: ts,
		Link:     &l,
	}
	return nil
}

// DeleteNetworkInterface :
func (n *NetworkInterfaceCtrl) DeleteNetworkInterface(name string) error {
	lattr := netlink.NewLinkAttrs()
	lattr.Name = name

	ipRoutes, err := netlink.RouteList(&netlink.Wireguard{LinkAttrs: lattr}, 0)
	if err != nil {
		return fmt.Errorf("failed to get IP routes list: %v", err)
	}
	for _, route := range ipRoutes {
		if err = netlink.RouteDel(&route); err != nil {
			return fmt.Errorf("failed to remove IP route: %v", err)
		}
	}

	if err := netlink.LinkDel(&netlink.Wireguard{LinkAttrs: lattr}); err != nil {
		return fmt.Errorf("failed to delete network device: %v", err)
	}

	return nil
}

// GetWgPrivateKey :
func (n *NetworkInterfaceCtrl) GetWgPrivateKey() *wgtypes.Key {
	return n.wgPrivateKey
}
