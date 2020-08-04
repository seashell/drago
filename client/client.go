package client

import (
	"context"
	"net"
	"reflect"
	"strconv"
	"time"

	api "github.com/seashell/drago/api"
	nic "github.com/seashell/drago/client/nic"
	"github.com/seashell/drago/client/storage"

	log "github.com/seashell/drago/pkg/log"
	netlink "github.com/vishvananda/netlink"
	wgtypes "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Client :
type Client struct {
	config *Config
	logger log.Logger

	niCtrl    *nic.NetworkInterfaceCtrl
	apiClient *api.Client
	stateDB   storage.StateRpository
}

// Config :
type Config struct {
	Enabled                  bool
	Servers                  []string
	Token                    string
	DataDir                  string
	InterfacesPrefix         string
	SyncInterval             time.Duration
	LinksPersistentKeepalive int
}

// New :
func New(conf *Config, log log.Logger) (*Client, error) {
	a, err := api.NewClient(&api.Config{
		Address: conf.Servers[0],
		Token:   conf.Token,
	})
	if err != nil {
		return nil, err
	}

	n, err := nic.NewCtrl(conf.InterfacesPrefix, conf.WireguardPath, log)
	if err != nil {
		return nil, err
	}

	s, err := storage.NewJsonStorage(conf.DataDir + "state.json")
	if err != nil {
		return nil, err
	}

	return &Client{
		config:    conf,
		niCtrl:    n,
		apiClient: a,
		stateDB:   s,
		logger:    log,
	}, nil
}

// Run :
func (c *Client) Run() {

	c.logger.Debugf("Applying local settings\n")

	ls, err := c.stateDB.GetHostSettings()
	if err != nil {
		c.logger.Warnf("Parsing error: %v\n", err)
	}
	if err := c.niCtrl.Update(c.fromApiSettingsToNic(ls)); err != nil {
		c.logger.Warnf("Interfaces update error: %v\n", err)
	} else {
		c.logger.Debugf("Finished applying local settings\n")
	}

	c.logger.Infof("Starting sychronization with servers every %v\n", c.config.SyncInterval)

	go func() {
		for {

			niState := []*api.WgInterfaceState{}

			for _, iface := range c.niCtrl.NetworkInterfaces {
				pk := iface.Settings.Wireguard.PrivateKey.PublicKey().String()
				niState = append(niState, &api.WgInterfaceState{
					Name:      iface.Settings.Alias,
					PublicKey: &pk,
				})
			}

			state := &api.HostState{
				Interfaces: niState,
			}

			desired, err := c.apiClient.Agent().SynchronizeSelf(context.Background(), state)

			if err == nil {

				current, err := c.stateDB.GetHostSettings()
				if err != nil {
					c.logger.Warnf("Error fetching host settings from local storage: %v\n", err)
				}

				if !reflect.DeepEqual(desired, current) {
					c.logger.Debugf("Applying remote settings\n")
					if err := c.niCtrl.Update(c.fromApiSettingsToNic(desired)); err != nil {
						c.logger.Warnf("Error applying network interface settings: %v\n", err)
					} else {
						if err := c.stateDB.PutHostSettings(desired); err != nil {
							c.logger.Warnf("Error persisting host settings to local storage: %v\n", err)
						}
						c.logger.Debugf("Host settings reconciliation successfully completed!")
					}
				}
			} else {
				c.logger.Warnf("Reconciliation error: %v\n", err)
			}

			time.Sleep(c.config.SyncInterval)

		}
	}()
}

func (c *Client) fromApiSettingsToNic(apiSettings *api.HostSettings) []nic.Settings {

	nicSettings := []nic.Settings{}

	for _, iface := range apiSettings.Interfaces {

		var wgPeers []wgtypes.PeerConfig

		for _, peer := range apiSettings.Peers {

			if *iface.Name == peer.Interface {
				var err error
				var pubkey wgtypes.Key
				if peer.PublicKey != nil {
					pubkey, err = wgtypes.ParseKey(*peer.PublicKey)
					if err != nil {
						c.logger.Warnf("Key parsing error: %v\n", err)
					}
				}

				var allowedIPs []net.IPNet
				for _, ip := range peer.AllowedIPs {
					_, allowedIP, err := net.ParseCIDR(ip)
					if err != nil {
						c.logger.Warnf("CIDR parsing error at %v: %v\n", err)
					}
					allowedIPs = append(allowedIPs, *allowedIP)
				}

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
					PublicKey:                   pubkey,
					AllowedIPs:                  allowedIPs,
					Endpoint:                    endpoint,
					PersistentKeepaliveInterval: persistentKeepalive,
				}
				wgPeers = append(wgPeers, wgPeer)
			}
		}

		var listenPort *int
		if iface.ListenPort != nil {
			lp, _ := strconv.Atoi(*iface.ListenPort)
			listenPort = &lp
		}

		wgConfig := &wgtypes.Config{
			PrivateKey:   c.niCtrl.GetWgPrivateKey(),
			ListenPort:   listenPort,
			ReplacePeers: true,
			Peers:        wgPeers,
		}

		addr, err := netlink.ParseAddr(*iface.Address)
		if err != nil {
			c.logger.Warnf("IP address parsing error at %v: %v\n", time.Now().Round(0), err)
		}

		nicSettings = append(nicSettings, nic.Settings{
			Alias:     iface.Name,
			Address:   addr,
			Wireguard: wgConfig,
		})

	}
	return nicSettings
}
