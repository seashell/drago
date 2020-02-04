package client

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"time"

	resty "github.com/go-resty/resty/v2"
	wg "github.com/squat/kilo/pkg/wireguard"
)

type client struct {
	config     ClientConfig
	httpClient *resty.Client
	key        []byte
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

	httpClient := resty.New()

	key, _ := wg.GenKey()

	if c.WgKey != "" {
		key = []byte(c.WgKey)
	}

	return &client{
		config:     c,
		key:        key,
		httpClient: httpClient,
	}, nil

}

func (c *client) PollConfigServer() *Node {

	url := "http://" + c.config.ServerAddr + "/api/v1/node"

	pk, _ := wg.PubKey(c.key)

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

func (c *client) Reconcile(n *Node) {

	const tmp = `
	[Interface]
	Address = {{ .Interface.Address }}
	ListenPort = {{ .Interface.ListenPort }}
	PrivateKey = {{ .Interface.PrivateKey }}
	
	{{with .Peers}}
	{{ range . }}
	[Peer]
	{{ if .Endpoint}}
	Endpoint = {{ .Endpoint }}
	{{else}}
	PublicKey = {{ .PublicKey }}
	AllowedIPs = {{ .AllowedIPs }}
	{{end}}
	PersistentKeepalive = {{ .PersistentKeepalive }}
	{{end}}
	{{end}}
	`

	t := template.Must(template.New("").Parse(tmp))

	err := t.Execute(os.Stdout, n)

	if err != nil {
		fmt.Println("Error executing template")
	}

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
