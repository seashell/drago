package host

//TODO improve separation of concerns between the API, the host, network management and wireguard settings

import (
	"time"
	"path"
	"fmt"
	"net"
	"errors"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"reflect"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"github.com/seashell/drago/client/networking"
	"github.com/seashell/drago/api"
)
const (
	defaultHostSettingsFile string = "settings.json"
)


type Config struct {
	DataDir string
	SyncInterval time.Duration
}

type Client struct {
	config *Config

	//networking controller client
	netClient *networking.Client

	//api client
	apiClient *api.Client
	hostSettings *api.Settings
	hostState *api.State

	//local host information
	hostSettingsFilePath string
}

func ReadLocalSettings(path string) (*api.Settings, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil,err
	}
	s := api.Settings{}
	err = json.Unmarshal([]byte(f), &s)
	if err != nil {
		return nil,err
	}
	return &s,nil
}

func New(c *Config) (*Client, error) {

	n,err := networking.New()
	if err != nil {
		fmt.Println("failed to create network controller client: ", err)
	}

	//Parse current Host settings from file	
	fp := path.Join(c.DataDir, defaultHostSettingsFile)
	hs,err := ReadLocalSettings(fp)
	if err != nil {
		fmt.Println("warning: failed to read local settings: ", err)
	}

	return &Client{
		config: c,
		netClient: n,
		apiClient: nil,
		hostSettings: hs,
		hostSettingsFilePath: fp,
	}, nil
}


func ParseIP(s string) (string, string, error) {
	ip, port, err := net.SplitHostPort(s)
	if err == nil {
		return ip, port, nil
	}

	ip2 := net.ParseIP(s)
	if ip2 == nil {
		return "","", errors.New("invalid IP")
	}

	return ip2.String(), "",nil
}


func (c *Client) SetAPIClient(a *api.Client) {
	c.apiClient = a
}

func (c *Client) parseLocalState() (*api.State, error) {
	//iterate interfaces and get name + publickey for each one
	var state api.State
	if c.hostSettings != nil {
		for _, iface := range c.hostSettings.Interfaces {
			state.Interfaces = append(state.Interfaces, api.InterfaceState{
				Name: iface.Name,
				PublicKey: c.netClient.GetWGPublicKey(),
			})
		}
	}
	return &state,nil
}

func (c *Client) persistRemoteSettings(rs *api.Settings) (error) {
	//save contents of rs to a file
	file, err := json.MarshalIndent(rs, "", " ")
	if err != nil {
		fmt.Println("failed to marshall JSON: ", err)
		return err
	}

	err = ioutil.WriteFile(c.hostSettingsFilePath, file, 0644)
	if err != nil {
		fmt.Println("failed to save JSON file: ", err)
		return err
	}

	c.hostSettings = rs
	return nil
}

func (c *Client) parseRemoteSettings(rs *api.Settings) ([]networking.IfaceConfig, error) {
	//TODO handle extra wireguard settings (dns, pre/post up ...)
	var ns []networking.IfaceConfig
	for _, iface := range rs.Interfaces { //iterate interfaces		

		//peers settings
		var peerList []wgtypes.PeerConfig
		for _,peer := range rs.Peers {//iterate peers

			if peer.Interface == iface.Name {//if peer belong to this interface
								
				var pubKey wgtypes.Key
				var err error
				if peer.PublicKey == "" {
					fmt.Println("warning: using dummy public key for peer")
					pubKey,err =  wgtypes.ParseKey(c.netClient.GetWGPublicKey())
					if err != nil {
						fmt.Println("warning: failed to parse public key: ",err)
					}
				} else {//if pub key not present yet, generate dummy base on local private key
					pubKey, err = wgtypes.ParseKey(peer.PublicKey)
					if err != nil { 
						fmt.Println("warning: failed to parse public key: ",err)
					}
				}

				var allowedIPs []net.IPNet
				for _,ip := range peer.AllowedIps {
					_,ipNet,err := net.ParseCIDR(ip)
					if err != nil {
						fmt.Println("warning: failed to parse IP CIDR: ",err)
					}
					allowedIPs = append(allowedIPs, *ipNet)
				}
				PersistentKeepalive := 20*time.Second
				peerConfig := wgtypes.PeerConfig{
					PublicKey: pubKey,
					Remove: false,
					UpdateOnly: false,
					PresharedKey: nil,
					ReplaceAllowedIPs: true,
					AllowedIPs: allowedIPs,
					PersistentKeepaliveInterval: &PersistentKeepalive,
				}
	
				if peer.Address != "" {
					//TODO include hostname parsing, remove port parsing
					host, port,_ := ParseIP(peer.Address)
		
					var ep int
					if port != "" {
						ep,_= strconv.Atoi(port)
					} else {
						ep,_ = strconv.Atoi(peer.Port)
					}
					endpoint := &net.UDPAddr{
						IP: net.ParseIP(host),
						Port: ep,
					}
					peerConfig.Endpoint = endpoint
				}				
				peerList = append(peerList, peerConfig)
			}			
		}	
		
		//interface settings
		listenPort,_ := strconv.Atoi(iface.ListenPort)
		Addr,err := networking.ToNetlinkAddr(iface.Address)
		if err != nil {
			fmt.Println("warning, failed convert address: ", err)
		}

		ns = append(ns, networking.IfaceConfig{
			Name: iface.Name,
			Address: &Addr,
			Wireguard: &wgtypes.Config{
				PrivateKey:   &c.netClient.WgPrivateKey,
				ListenPort:  &listenPort,
				Peers: peerList,
				ReplacePeers: true,
			},
		})

	}

	return ns,nil
}


func (c *Client) Start() {

	fmt.Println("Syncing with remote every",c.config.SyncInterval)
	go func() {
		for {
			//Update local state from system, 
			ls,err := c.parseLocalState()
			if err != nil {
				fmt.Println("warning, failed to parse local state from file: ", err)
			}
			//and sync with remote
			rs,err := c.apiClient.Hosts().PostSelfSync(ls)
			if err != nil {
				fmt.Println("warning, failed to sync with remote: ", err)
			}

			if rs != nil {// if remote state is not empty
				if ! reflect.DeepEqual(rs,c.hostSettings) {//if settings changed
					fmt.Println("Starting host settings update...")
					// persist new settings
					err := c.persistRemoteSettings(rs)
					if err != nil {
						fmt.Println("warning, failed to persist remote settigns locally: ", err)
					} else {
						// parse new settings into appropriate type
						ns,err := c.parseRemoteSettings(rs)
						if err != nil {
							fmt.Println("warning, failed to persist remote settigns locally: ", err)
						} else {
							// and apply them to the networking interfaces
							c.netClient.Apply(ns)
							if err != nil {
								fmt.Println("warning, failed to apply changes to host: ", err)
							}
							fmt.Println("Update done")
						}
					}					
				}	
			}
			time.Sleep(c.config.SyncInterval)
		}
	}()
}