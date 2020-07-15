package client

import (
	"net"
	"reflect"
	"strconv"
	"time"

	"github.com/seashell/drago/api"
	"github.com/seashell/drago/client/nic"
	"github.com/seashell/drago/client/state"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"github.com/seashell/drago/pkg/logger"
)

type Client struct {
	config Config

	niCtrl    *nic.NetworkInterfaceCtrl
	apiClient *api.Client

	stateDB state.StateDB

	log      logger.Logger
}

type Config struct {
	Enabled      		bool
	Servers      		[]string
	Token        		string
	DataDir      		string
	InterfacesPrefix	string
	SyncInterval 		time.Duration
	LinksPersistentKeepalive	int
}

func New(c Config, l logger.Logger) (*Client, error) {
	a, err := api.NewClient(&api.Config{
		Address: c.Servers[0], //TODO: add support for multiple API addresses
		Token:   c.Token,
	})
	if err != nil {
		return nil, err
	}

	n, err := nic.NewCtrl(c.InterfacesPrefix)
	if err != nil {
		return nil, err
	}

	s, err := state.NewFileDB(c.DataDir)
	if err != nil {
		return nil, err
	}

	return &Client{
		config:    	c,
		niCtrl:    	n,
		apiClient: 	a,
		stateDB:   	s,
		log:		l,
	}, nil
}

func (c *Client) Run() {
	c.log.Debugf("Applying local settings at %v\n", time.Now().Round(0))
	ls, err := c.stateDB.GetHostSettings()
	if err != nil {
		c.log.Warnf("Parsing error at %v: %v\n", time.Now().Round(0), err)
	}
	if err := c.niCtrl.Update(c.fromApiSettingsToNic(ls)); err != nil {
		c.log.Warnf("Interfaces update error at %v: %v\n", time.Now().Round(0), err)
	} else {
		c.log.Debugf("Finished applying local settings at %v\n", time.Now().Round(0))
	}

	c.log.Infof("Starting sychronization with servers every %v\n", c.config.SyncInterval)
	go func() {
		for {
			// Parse current host network interfaces state
			niState := []api.NetworkInterfaceState{}
			for _, iface := range c.niCtrl.NetworkInterfaces {
				niState = append(niState, api.NetworkInterfaceState{
					Name:        *iface.Settings.Alias,
					WgPublicKey: iface.Settings.Wireguard.PrivateKey.PublicKey().String(),
				})
			}
			// Submit current network interfaces state and get target remote settings
			ts, err := NewSynchronizationEndpoint(c).SynchronizeSelf(&api.HostState{NetworkInterfaces: niState})
			if err != nil {
				c.log.Warnf("Server synchronization error at %v: %v\n", time.Now().Round(0), err)
			} else if ts != nil {
				ls, err := c.stateDB.GetHostSettings()
				if err != nil {
					c.log.Warnf("Database error at %v: %v\n", time.Now().Round(0), err)
				}
				//If target remote settings != local settings, apply remote settings
				if !reflect.DeepEqual(ts, ls) {
					c.log.Debugf("Applying remote settings at %v\n", time.Now().Round(0))
					if err := c.niCtrl.Update(c.fromApiSettingsToNic(ts)); err != nil {
						//If not sucessful,  do not persist remote settings
						c.log.Warnf("Interfaces update error at %v: %v\n", time.Now().Round(0), err)
					} else {
						//If successful, update target local settings with target remote settings
						if err := c.stateDB.PutHostSettings(ts); err != nil {
							c.log.Warnf("Database error at %v: %v\n", time.Now().Round(0), err)
						}
						c.log.Debugf("Finished applying remote settings at %v\n", time.Now().Round(0))
					}
				}
			}

			time.Sleep(c.config.SyncInterval)
		}
	}()

}

func (c *Client) fromApiSettingsToNic(as *api.HostSettings) []nic.Settings {
	ts := []nic.Settings{}
	for _, ni := range as.NetworkInterfaces {
		//Parse WG settings

		// peers
		var wgPeers []wgtypes.PeerConfig
		for _, peer := range as.WireguardPeers {
			if *ni.Name == *peer.Interface {
				//Key
				var pub wgtypes.Key
				var err error
				if peer.PublicKey != nil {
					pub, err = wgtypes.ParseKey(*peer.PublicKey)
					if err != nil {
						c.log.Warnf("Key parsing error at %v: %v\n", time.Now().Round(0), err)
					}
				}

				//AllowedIPs
				var allowedIPs []net.IPNet
				for _, ip := range peer.AllowedIps {
					_, allowedIP, err := net.ParseCIDR(ip)
					if err != nil {
						c.log.Warnf("CIDR parsing error at %v: %v\n", time.Now().Round(0), err)
					}
					allowedIPs = append(allowedIPs, *allowedIP)
				}

				//PersistentKeepalive
				var persistentKeepalive *time.Duration
				if peer.PersistentKeepalive != nil {
					t := time.Duration(*peer.PersistentKeepalive) * time.Second
					persistentKeepalive = &t
				} else {
					if c.config.LinksPersistentKeepalive != 0 {
						t := time.Duration(c.config.LinksPersistentKeepalive) * time.Second
						persistentKeepalive = &t
					}					
				}

				//Endpoint
				var endpoint *net.UDPAddr
				if peer.Address != nil {
					p, _ := strconv.Atoi(*peer.Port)
					endpoint = &net.UDPAddr{
						IP:   net.ParseIP(*peer.Address),
						Port: p,
					}
				}

				wgPeer := wgtypes.PeerConfig{
					Remove:                      false,
					UpdateOnly:                  false,
					ReplaceAllowedIPs:           true,
					PresharedKey:                nil,
					PublicKey:                   pub,
					AllowedIPs:                  allowedIPs,
					Endpoint:                    endpoint,
					PersistentKeepaliveInterval: persistentKeepalive,
				}
				wgPeers = append(wgPeers, wgPeer)
			}
		}

		//ListenPort
		var listenPort *int
		if ni.ListenPort != nil {
			lp, _ := strconv.Atoi(*ni.ListenPort)
			listenPort = &lp
		}

		wgConfig := &wgtypes.Config{
			PrivateKey:   c.niCtrl.GetWgPrivateKey(),
			ListenPort:   listenPort,
			ReplacePeers: true,
			Peers:        wgPeers,
		}

		//Parse link device settings
		//Address
		addr, err := netlink.ParseAddr(*ni.Address)
		if err != nil {
			c.log.Warnf("IP address parsing error at %v: %v\n", time.Now().Round(0), err)
		}
		ts = append(ts, nic.Settings{ 
			Alias:	   ni.Name,
			Address:   addr,
			Wireguard: wgConfig,
		})
	}
	return ts
}
