package nic

import (
	"strings"

	netlink "github.com/vishvananda/netlink"
)

func getLinkIndex(l string) (int, error) {
	link, err := netlink.LinkByName(l)
	if err != nil {
		return 0, err
	}
	return link.Attrs().Index, nil
}

func setLinkAlias(l string, alias string) error {

	link, err := netlink.LinkByName(l)
	if err != nil {
		return err
	}

	err = netlink.LinkSetAlias(link, alias)
	if err != nil {
		return err
	}

	return nil
}

func enableLink(l string) error {

	link, err := netlink.LinkByName(l)
	if err != nil {
		return err
	}

	err = netlink.LinkSetUp(link)
	if err != nil {
		return err
	}

	return nil
}

func deleteLinkAndRoutesWithPrefix(p string) error {

	links, err := listLinksWithPrefix(p)
	if err != nil {
		return err
	}

	for _, l := range links {
		err := deleteLinkAndRoutes(l.Attrs().Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func listLinksWithPrefix(p string) ([]netlink.Link, error) {

	out := []netlink.Link{}

	links, _ := netlink.LinkList()
	for _, l := range links {
		if strings.HasPrefix(l.Attrs().Name, p) {
			out = append(out, l)
		}
	}

	return out, nil
}

func deleteLinkAndRoutes(l string) error {

	attrs := netlink.NewLinkAttrs()
	attrs.Name = l

	routes, err := listLinkRoutes(l)
	if err != nil {
		return err
	}

	for _, r := range routes {
		if err = netlink.RouteDel(&r); err != nil {
			return err
		}
	}

	if err := netlink.LinkDel(&netlink.Wireguard{LinkAttrs: attrs}); err != nil {
		return err
	}

	return nil
}

func listLinkRoutes(l string) ([]netlink.Route, error) {

	attrs := netlink.NewLinkAttrs()
	attrs.Name = l

	routes, err := netlink.RouteList(&netlink.Wireguard{LinkAttrs: attrs}, 0)
	if err != nil {
		return nil, err
	}

	return routes, nil
}

func setLinkAddress(l string, cidr string) error {

	link, err := netlink.LinkByName(l)
	if err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(cidr)
	if err != nil {
		return err
	}

	err = netlink.AddrAdd(link, addr)
	if err != nil {
		return err
	}

	return err
}
