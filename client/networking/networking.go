package networking

import (
	"fmt"
	"os"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type IfaceConfig struct {
	Name string
	Address *netlink.Addr
	Wireguard *wgtypes.Config
}

type Iface struct {
	config *IfaceConfig
	link netlink.Link
}

type Client struct {
	ifaces map[string]*Iface	
	wgCtrl *wgctrl.Client

	WgPrivateKey wgtypes.Key
}


func New() (*Client,error) {

	wg, err := wgctrl.New()
	if err != nil {
		fmt.Println("failed to create wireguard client: ", err)
	}

	priv, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		fmt.Println("failed to generate wireguard private key:", err) 
		panic(err)
	}

	return &Client{
		ifaces: make(map[string]*Iface),
		wgCtrl: wg,
		WgPrivateKey: priv,
	}, nil
}

func (c *Client) GetWGPublicKey() (string) {
	return c.WgPrivateKey.PublicKey().String()
}

func (c *Client) Apply(ifacesConfigs []IfaceConfig) (error) {
	for _, iface := range ifacesConfigs {

		if c.ifaces[iface.Name] == nil { //if this is a new interface, register it
			if err := c.registerIface(iface.Name); err != nil {
				return err
			}
		}

		if c.ifaces[iface.Name].config == nil {
			if err := netlink.AddrReplace(c.ifaces[iface.Name].link, iface.Address); err != nil {
				fmt.Println("failed to add address:", err)
				return err
			}	
		} else if ! iface.Address.Equal(*c.ifaces[iface.Name].config.Address) { //if different than the new one, add address and delete previous
			if err := netlink.AddrAdd(c.ifaces[iface.Name].link, iface.Address); err != nil {
				fmt.Println("failed to change adress:", err)
				return err
			}
		}		

		addresses,_ := netlink.AddrList(c.ifaces[iface.Name].link, 0) //get current iface address	
		for _,addr := range addresses {
			if ! addr.Equal(*iface.Address) { //if different than the new one, add address and delete previous
				if err := netlink.AddrDel(c.ifaces[iface.Name].link, &addr); err != nil {
					fmt.Println("failed to delete address:", err)
					return err
				}
			}
		}

		//wireguard
		if err := c.wgCtrl.ConfigureDevice(iface.Name, *iface.Wireguard); err != nil {
			if os.IsNotExist(err) {
				fmt.Println("failed to configure wireguard device: ", err)
				return err
			} else {
				fmt.Println("Unknown device configuration error: ", err)
				return err
			}
		}
	
		c.ifaces[iface.Name].config = &iface
	}
	// TODO: remove dangling ifaces
	return nil
}

func (c *Client) registerIface(name string) (error) {
	iface,err := c.linkIface(name)
	if err != nil {
		fmt.Println("failed to link interface: ", err)
		return err
	}

	c.ifaces[name] = iface
	return nil
}

func (c *Client) linkIface(name string) (*Iface, error) {
	//try to add new interface at the host
	la := netlink.NewLinkAttrs()
	la.Name = name
	err := netlink.LinkAdd(&netlink.Wireguard{LinkAttrs: la})
	if err != nil  {
		//in case it fails because interface already exists, we just link it by name
		fmt.Println("Warning: failed to create new link device: ", err)
	}
	
	//link interface by name
	l, err := netlink.LinkByName(name)
	if err != nil {
		fmt.Println("failed to get link device by name: ", err)
		return nil,err
	}
	return &Iface{
		config: nil,
		link: l,
	},nil

}


func ToNetlinkAddr (addr string) (netlink.Addr, error) {
	//parse IP address into netlink.Addr format
	a,err := netlink.ParseAddr(addr) 
	if err != nil {
		fmt.Println("failed to parse IP address:", err)
		return netlink.Addr{}, err
	}
	
	return *a, nil
}