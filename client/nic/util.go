package nic

import (
	"strings"

	netlink "github.com/vishvananda/netlink"
)

func linkNameByAlias(s string) (string, error) {
	link, err := netlink.LinkByAlias(s)
	if err != nil {
		return "", err
	}
	return link.Attrs().Name, nil

}

func deleteLinkAndRoutesByName(s string) error {

	link, err := netlink.LinkByName(s)
	if err != nil {
		return err
	}

	deleteLinkAndRoutes(link)

	return nil
}

func deleteLinksAndRoutesByAlias(s string) error {

	links, err := linksByAlias(s)
	if err != nil {
		return err
	}

	for _, l := range links {
		err := deleteLinkAndRoutes(l)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteLinksAndRoutesByPrefix(s string) error {

	links, err := linksByPrefix(s)
	if err != nil {
		return err
	}

	for _, l := range links {
		err := deleteLinkAndRoutes(l)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteLinkAndRoutes(link netlink.Link) error {

	// List routes
	routes, err := netlink.RouteList(&netlink.Wireguard{LinkAttrs: *link.Attrs()}, 0)
	if err != nil {
		return err
	}

	// Delete routes
	for _, r := range routes {
		if err = netlink.RouteDel(&r); err != nil {
			return err
		}
	}

	// Delete link
	if err := netlink.LinkDel(&netlink.Wireguard{LinkAttrs: *link.Attrs()}); err != nil {
		return err
	}

	return nil
}

func linksByPrefix(p string) ([]netlink.Link, error) {

	out := []netlink.Link{}

	links, _ := netlink.LinkList()
	for _, l := range links {
		if strings.HasPrefix(l.Attrs().Name, p) {
			out = append(out, l)
		}
	}

	return out, nil
}

func linksByAlias(s string) ([]netlink.Link, error) {

	out := []netlink.Link{}

	links, _ := netlink.LinkList()
	for _, l := range links {
		if l.Attrs().Alias == s {
			out = append(out, l)
		}
	}

	return out, nil
}

func setLinkAddress(link netlink.Link, cidr string) error {

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
