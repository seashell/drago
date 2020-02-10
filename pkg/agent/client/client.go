package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	"net"

	resty "github.com/go-resty/resty/v2"
	wg "github.com/squat/kilo/pkg/wireguard"
	ip "github.com/squat/kilo/pkg/iproute"
)

type client struct {
	key        string
	config     ClientConfig
	httpClient *resty.Client
	wgConf     *wg.Conf
	wgIface    int
	node	   *Node
}

type ClientConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	SyncInterval int    `mapstructure:"sync_interval"`
	ServerAddr   string `mapstructure:"server_addr"`
	WgKey        string `mapstructure:"wg_key"`
	Jwt          string `mapstructure:"jwt"`
	DataDir		 string `mapstructure:"data_dir"`
}

type Client interface {
	Run()
}

func New(c ClientConfig) (*client, error) {

	wgIface, isNew, err := wg.New("wg0")
	if err != nil {
		return nil, fmt.Errorf("%v",err)
	}

	if isNew {
		fmt.Println("New Wireguard interface create") //???
	}

	dat, err := ioutil.ReadFile(c.DataDir + "/wg0.conf")
	if err != nil {
		fmt.Println("Error reading Wireguard config file: ",err)
	}

	n := &Node{}
	node, err := ioutil.ReadFile(c.DataDir + "/wg0.json")
	if err != nil {
		fmt.Println("Error reading node JSON file: ",err)
	} else {
		err = json.Unmarshal(node, n)
		if err != nil {
			fmt.Println("Error unmarshalling node JSON: ",err)
		}
	}


	// If a key already exists, use it. Otherwise generate a new one.
	key := c.WgKey
	if key == "" {
		k, _ := wg.GenKey()
		key = string(k)
	}

	return &client{
		config:     c,
		key:        key,
		wgIface:    wgIface,
		wgConf:     wg.Parse(dat),
		httpClient: resty.New(),
		node:		n,
	}, nil

}


func (c *client) PollConfigServer() (*Node, error) {

	url := "http://" + c.config.ServerAddr + "/api/v1/node"

	// Update the server with our current public key, to allow for key rotation
	pk, _ := wg.PubKey([]byte(c.key))
	body := map[string]string{
		"publicKey": string(pk),
	}

	resp, err := c.httpClient.R().
		SetHeader("Authorization", "Bearer "+c.config.Jwt).
		SetBody(body).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("Error polling server: %v",err)
	}

	n := &Node{}
	err = json.Unmarshal(resp.Body(), n)
	n.Interface.PrivateKey = c.config.WgKey

	return n, nil
}

func (c *client) Reconcile(n *Node) error {

	if n.Interface.Address != c.node.Interface.Address {

		path := c.config.DataDir+"/wg0.json"
		buf, err := json.Marshal(n)
		if err != nil {
			return fmt.Errorf("Error formatting JSON: %v",err)
		}
		if err := ioutil.WriteFile(path, buf, 0600); err != nil {
			return fmt.Errorf("Error writing node JSON to file: %v",err)
		}

		_, ifaceAddr, err := net.ParseCIDR(n.Interface.Address)
		if err != nil {
			return fmt.Errorf("Error parsing interface CIDR: %v",err)
		}

		if err = ip.SetAddress(c.wgIface, ifaceAddr); err != nil{
			return fmt.Errorf("Error setting interface IP address: %v",err)
		}	

		c.node = n
	}

	conf, err := templateToString(tmpl, n)
	if err != nil {
		return fmt.Errorf("Error converting Wireguard config to string: %v",err)
	}

	newConf := wg.Parse([]byte(conf))
	if equal := newConf.Equal(c.wgConf); !equal {

		path := c.config.DataDir+"/wg0.conf"
		if err := templateToFile(tmpl, path, n); err != nil {
			return fmt.Errorf("Error writing Wireguard config to file: %v",err)
		}
		
		if err = wg.SetConf("wg0", path); err != nil{
			return fmt.Errorf("Error applying Wiregard conf: %v",err)
		}

		c.wgConf = newConf
	}
	
	if err = ip.Set(c.wgIface, true); err != nil {
		return fmt.Errorf("Error bringing interface up: %v",err)
	}

	return nil

}

func (c *client) Run() {
	go func() {
		for {
			time.Sleep(time.Duration(c.config.SyncInterval) * time.Second)
			n, err := c.PollConfigServer()
			if err != nil {
				fmt.Println(err.Error())
			}
			if err := c.Reconcile(n); err != nil{
				fmt.Println(err.Error())
			}
		}
	}()
}
