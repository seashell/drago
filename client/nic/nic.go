package nic

import (
	"fmt"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Settings struct {
	Name      *string
	Address   *netlink.Addr
	Wireguard *wgtypes.Config
}

type NetworkInterface struct {
	Settings *Settings
	Link     *netlink.Link
}

type NetworkInterfaceCtrl struct {
	NetworkInterfaces map[string]*NetworkInterface

	wgController *wgctrl.Client
	wgPrivateKey *wgtypes.Key
}

func NewCtrl() (*NetworkInterfaceCtrl, error) {

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
	}, nil
}

func (n *NetworkInterfaceCtrl) Update(ts []Settings) error {
	if err := n.resetWgNetworkInterfaces(); err != nil {
		return err
	}

	for _, s := range ts {
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
			if err := n.DeleteNetworkInterface(ni.Attrs().Name); err != nil {
				fmt.Println("Warning: failed to delete network interface: ", err)
			}
			delete(n.NetworkInterfaces, ni.Attrs().Name)
		}

	}
	return nil
}

func (n *NetworkInterfaceCtrl) ConfigureNetworkInterface(ts *Settings) error {
	// register new interface
	lattr := netlink.NewLinkAttrs()
	lattr.Name = *ts.Name
	if err := netlink.LinkAdd(&netlink.Wireguard{LinkAttrs: lattr}); err != nil {
		fmt.Println("failed to create new network device: ", err)
		return err
	}

	l, err := netlink.LinkByName(*ts.Name)
	if err != nil {
		fmt.Println("failed to get network device by name: ", err)
		return err
	}

	// apply wireguard config
	if err := n.wgController.ConfigureDevice(*ts.Name, *ts.Wireguard); err != nil {
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

	n.NetworkInterfaces[*ts.Name] = &NetworkInterface{
		Settings: ts,
		Link:     &l,
	}
	return nil
}

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

func (n *NetworkInterfaceCtrl) GetWgPrivateKey() *wgtypes.Key {
	return n.wgPrivateKey
}
