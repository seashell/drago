package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"


	"github.com/vishvananda/netlink"
	resty "github.com/go-resty/resty/v2"
	wg "github.com/squat/kilo/pkg/wireguard"
)

type client struct {
	config     ClientConfig
	httpClient *resty.Client
	wgIface    int
	wgConf     *wg.Conf
	host       *Host
}

type ClientConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	DataDir      string `mapstructure:"data_dir"`
	Iface        string `mapstructure:"iface"`
	ServerAddr   string `mapstructure:"server_addr"`
	WgKey        string `mapstructure:"wg_key"`
	Jwt          string `mapstructure:"jwt"`
	SyncInterval int    `mapstructure:"sync_interval"`
}

func New(c ClientConfig) (*client, error) {

	wgIface, isNew, err := wg.New(c.Iface)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if isNew {
		fmt.Println("New Wireguard interface created")
	}

	err = os.MkdirAll(c.DataDir, 0600)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	dat, err := ioutil.ReadFile(c.DataDir + "/" + c.Iface + ".conf")
	if err != nil {
		fmt.Println("Error reading Wireguard config file: ", err)
	}

	h := &Host{}
	host, err := ioutil.ReadFile(c.DataDir + "/" + c.Iface + ".json")
	if err != nil {
		fmt.Println("Error reading host JSON file: ", err)
	} else {
		err = json.Unmarshal(host, h)
		if err != nil {
			fmt.Println("Error unmarshalling host JSON: ", err)
		}
	}

	// If a key already exists, use it. Otherwise generate a new one.
	privateKey := c.WgKey
	if privateKey == "" {
		k, _ := wg.GenKey()
		privateKey = string(k)
	}

	h.Keys.PrivateKey = privateKey

	publicKey, _ := wg.PubKey([]byte(privateKey))
	h.Keys.PublicKey = string(publicKey)

	return &client{
		config:     c,
		wgIface:    wgIface,
		wgConf:     wg.Parse(dat),
		httpClient: resty.New(),
		host:       h,
	}, nil

}

func (c *client) PollConfigServer() (*Host, error) {

	url := "http://" + c.config.ServerAddr + "/api/v1/hosts/self/settings"

	// Update the server with our current public key, to allow for key rotation
	body := map[string]string{
		"publicKey": c.host.Keys.PublicKey,
	}

	resp, err := c.httpClient.R().
		SetHeader("Authorization", "Bearer "+c.config.Jwt).
		SetBody(body).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("Error polling server: %v", err)
	}

	if resp.StatusCode() > 299 {
		return nil, fmt.Errorf("Error polling server: %v %v", resp.StatusCode(), resp)
	}

	h := &Host{}
	err = json.Unmarshal(resp.Body(), h)

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling Host: %v", err)
	}

	h.Keys = c.host.Keys
	h.AdvertiseAddr = c.host.AdvertiseAddr
	h.Jwt = c.host.Jwt

	return h, nil
}

func (c *client) Reconcile(h *Host) error {

	lo, err := netlink.LinkByIndex(c.wgIface)
	if err != nil {
		return fmt.Errorf("Error getting link to interface: %v", err)
	}

	if h.Interface.Address != c.host.Interface.Address {

		path := c.config.DataDir + "/" + c.config.Iface + ".json"
		buf, err := json.Marshal(h)
		if err != nil {
			return fmt.Errorf("Error formatting JSON: %v", err)
		}

		if err := ioutil.WriteFile(path, buf, 0600); err != nil {
			return fmt.Errorf("Error writing host JSON to file: %v", err)
		}

		addr, err := netlink.ParseAddr(h.Interface.Address)
		if err != nil {
			return fmt.Errorf("Error parsing IP address: %v", err)
		}

		if err = netlink.AddrAdd(lo, addr); err != nil {
			return fmt.Errorf("Error setting interface IP address: %v", err)
		}

		c.host = h
	}

	conf, err := templateToString(tmpl, h)
	if err != nil {
		return fmt.Errorf("Error converting Wireguard config to string: %v", err)
	}

	newConf := wg.Parse([]byte(conf))
	if equal := newConf.Equal(c.wgConf); !equal {

		path := c.config.DataDir + "/" + c.config.Iface + ".conf"
		if err := templateToFile(tmpl, path, h); err != nil {
			return fmt.Errorf("Error writing Wireguard config to file: %v", err)
		}

		if err = wg.SetConf(c.config.Iface, path); err != nil {
			return fmt.Errorf("Error applying Wiregard conf: %v", err)
		}

		c.wgConf = newConf
	}

	if err = netlink.LinkSetUp(lo); err != nil {
		return fmt.Errorf("Error bringing interface up: %v", err)
	}

	return nil

}

func (c *client) Run() {
	go func() {
		for {
			time.Sleep(time.Duration(c.config.SyncInterval) * time.Second)
			h, err := c.PollConfigServer()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				if err := c.Reconcile(h); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}()
}
