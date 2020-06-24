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
		//Bring Iface downif it exists
		l, _ := netlink.LinkByName(iface.Name)
		if l != nil {
			if err := c.bringIfaceDown(l); err != nil {
				return err
			}
			c.ifaces[iface.Name] = nil //clear iface
		}
		
		//TODO: Preup hook		
		//register interface
		if c.ifaces[iface.Name] == nil {
			if err := c.registerIface(iface.Name); err != nil {
				return err
			}
		}
		
		//Configure wireguard
		if err := c.wgCtrl.ConfigureDevice(iface.Name, *iface.Wireguard); err != nil {
			if os.IsNotExist(err) {
				fmt.Println("failed to configure wireguard device: ", err)
				return err
			} else {
				fmt.Println("Unknown device configuration error: ", err)
				return err
			}
		}
		
		//Add address 
		if err := netlink.AddrAdd(c.ifaces[iface.Name].link, iface.Address); err != nil {
			fmt.Println("failed to add address:", err)
			return err
		}	
		
		//TODO: Set MTU

		//bring interface Up
		if err := netlink.LinkSetUp(c.ifaces[iface.Name].link); err != nil {
			fmt.Println("failed to set link up:", err)
			return err
		}

		//TODO: Set DNS
		//TODO: set ip routes
		//TODO: Postup hook
	
		c.ifaces[iface.Name].config = &iface
	}
	
	//Delete dangling interfaces
	for existingIfaceName,existingIface := range c.ifaces {
		dangling := true
		for _,targetIface := range ifacesConfigs {
			if targetIface.Name == existingIfaceName {
				dangling = false
			}
		}
		if dangling {
			if err := c.bringIfaceDown(existingIface.link); err != nil {
				fmt.Println("failed to delete interface")
				return err
			}
			delete(c.ifaces, existingIfaceName)
		}
	}

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

func (c *Client) bringIfaceDown(l netlink.Link) (error) {

		//TODO: PreDown hook
		//delete interface
		//TODO: unset DNS
		//TODO: remove firewall
		//TODO: clean ip routes
		err := netlink.LinkDel(l)
		if err != nil  {
			//in case it fails because interface already exists, we just link it by name
			fmt.Println("Warning: failed to delte link device: ", err)
		}
		//TODO: PostDown hook
	return nil
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