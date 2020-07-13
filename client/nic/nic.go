package nic

import (
	"fmt"
	"math/rand"
	"encoding/hex"
	"regexp"	

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Settings :
type Settings struct {
	Name      	string
	Alias		*string	
	Address   	*netlink.Addr
	Wireguard	*wgtypes.Config
}

// NetworkInterface :
type NetworkInterface struct {
	Settings *Settings
	Link     *netlink.Link
}

// NetworkInterfaceCtrl :
type NetworkInterfaceCtrl struct {
	NetworkInterfaces map[string]*NetworkInterface

	namePrefix	 string
	wgController *wgctrl.Client
	wgPrivateKey *wgtypes.Key
}

// NewCtrl :
func NewCtrl(namePrefix string) (*NetworkInterfaceCtrl, error) {

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
		wgController:   wg,
		wgPrivateKey:   &pk,
		namePrefix:		namePrefix,
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
		s.Name = n.namePrefix+r
		if err := n.ConfigureNetworkInterface(&s); err != nil {
			return err
		}
	}
	return nil
}

func (n *NetworkInterfaceCtrl) resetWgNetworkInterfaces() error {
	niList, _ := netlink.LinkList()
	for _, ni := range niList {
		if ni.Type() == "wireguard" {
			//match device alias with prefix provided by n.namePrefix 
			matched, err := regexp.MatchString(n.namePrefix+`.*`, ni.Attrs().Name)
			if err != nil {
				fmt.Println("Warning: failed to match interface name: ", err)
			}
		
			if matched {
				if err := n.DeleteNetworkInterface(ni.Attrs().Name); err != nil {
					fmt.Println("Warning: failed to delete network interface: ", err)
				}
				delete(n.NetworkInterfaces, ni.Attrs().Alias)
			}

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

	if err := netlink.LinkAdd(&netlink.Wireguard{LinkAttrs: lattr}); err != nil {
		fmt.Println("failed to create new network device: ", err)
		return err
	}

	l, err := netlink.LinkByName(ts.Name)
	if err != nil {
		fmt.Println("failed to get network device by name: ", err)
		return err
	}

	if err := netlink.LinkSetAlias(l, lattr.Alias); err != nil {
		fmt.Println("failed to set link alias: ", err)
		return err
	}


	// apply wireguard config
	if err := n.wgController.ConfigureDevice(ts.Name, *ts.Wireguard); err != nil {
		if err != nil {
			fmt.Println("Unknown device configuration error: ", err)
			return err
		}
	}

	// apply new settings
	if err := netlink.AddrAdd(l, ts.Address); err != nil {
		fmt.Println("failed to add IP address:", err)
		return err
	}

	if err := netlink.LinkSetUp(l); err != nil {
		fmt.Println("failed to set network device up:", err)
		return err
	}

	for _, peer := range ts.Wireguard.Peers {
		for _, allowedIP := range peer.AllowedIPs {
			if err = netlink.RouteAdd(&netlink.Route{
				LinkIndex: l.Attrs().Index,
				Dst:       &allowedIP,
			}); err != nil {
				fmt.Println("failed to add IP route:", err)
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
		fmt.Println("warning: failed to get IP routes list:", err)
	}
	for _, route := range ipRoutes {
		if err = netlink.RouteDel(&route); err != nil {
			fmt.Println("warning: failed to remove IP route:", err)
		}
	}

	if err := netlink.LinkDel(&netlink.Wireguard{LinkAttrs: lattr}); err != nil {
		fmt.Println("Warning: failed to delete network device: ", err)
	}

	return nil
}

// GetWgPrivateKey :
func (n *NetworkInterfaceCtrl) GetWgPrivateKey() *wgtypes.Key {
	return n.wgPrivateKey
}
