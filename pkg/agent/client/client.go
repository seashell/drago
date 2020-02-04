package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"time"

	resty "github.com/go-resty/resty/v2"
	wg "github.com/squat/kilo/pkg/wireguard"
)

type client struct {
	key        string
	config     ClientConfig
	httpClient *resty.Client
	wgConf     *wg.Conf
	wgIface    int
}

type ClientConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	SyncInterval int    `mapstructure:"sync_interval"`
	ServerAddr   string `mapstructure:"server_addr"`
	WgKey        string `mapstructure:"wg_key"`
	Jwt          string `mapstructure:"jwt"`
}

type Client interface {
	Run()
}

func New(c ClientConfig) (*client, error) {

	wgIface, isNew, err := wg.New("wg0")
	if err != nil {
		fmt.Println("Error initializing Wireguard interface")
	}

	if isNew {
		fmt.Println("New Wireguard interface create")
	}

	dat, err := ioutil.ReadFile("./wg0.conf")
	if err != nil {
		fmt.Println("Error reading Wireguard configuration file")
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
	}, nil

}

func (c *client) PollConfigServer() *Node {

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
		fmt.Println(err)
	}

	n := &Node{}
	err = json.Unmarshal(resp.Body(), n)

	n.Interface.PrivateKey = c.config.WgKey

	return n
}

func (c *client) Reconcile(n *Node) error {

	path := "./wg0.conf"

	if err := templateToFile(tmpl, path, n); err != nil {
		fmt.Println("Error writing Wireguard config to file")
	}

	conf, err := templateToString(tmpl, n)
	if err != nil {
		fmt.Println("Error converting Wireguard config to string")
	}

	newConf := wg.Parse([]byte(conf))

	if newConf.Equal(c.wgConf) {
		fmt.Println("Configuration is still the same. No need for reconciliation.")
		return nil
	}

	fmt.Println("Reconciling Wireguard configuration state.")
	c.wgConf = newConf

	wg.ShowConf("wg0")

	return nil

}

func (c *client) Run() {

	go func() {
		for {
			time.Sleep(time.Duration(c.config.SyncInterval) * time.Second)
			n := c.PollConfigServer()
			c.Reconcile(n)
		}
	}()
}

func templateToFile(tmpl string, path string, ctx interface{}) error {

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()
	t := template.Must(template.New("").Parse(tmpl))

	err = t.Execute(f, ctx)
	if err != nil {
		return err
	}

	return nil
}

func templateToString(tmpl string, ctx interface{}) (string, error) {

	t := template.Must(template.New("").Parse(tmpl))

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, ctx); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
