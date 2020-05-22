package command

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type StorageStanza struct {
	Type string `hcl:"type,label"`
	Path string `hcl:"path,optional"`
}

type ServerStanza struct {
	Enabled bool           `hcl:"enabled,optional" env:"SERVER"`
	DataDir string         `hcl:"data_dir,optional"`
	Storage *StorageStanza `hcl:"storage,block"`
}

type ClientStanza struct {
	Enabled bool     `hcl:"enabled,optional" env:"CLIENT"`
	Servers []string `hcl:"servers,optional"`
	DataDir string   `hcl:"data_dir,optional"`
}

type VaultStanza struct {
	Enabled bool    `hcl:"enabled,optional"`
	Address *string `hcl:"address,optional"`
}

type DragoConfig struct {
	UI       bool          `hcl:"ui,optional"`
	DataDir  string        `hcl:"data_dir,optional"`
	BindAddr string        `hcl:"bind_addr,optional"`
	Server   *ServerStanza `hcl:"server,block"`
	Client   *ClientStanza `hcl:"client,block"`
	Vault    *VaultStanza  `hcl:"vault,block"`
}

func LoadConfigFromFile(path string) DragoConfig {

	c := DragoConfig{}

	err := hclsimple.DecodeFile(path, nil, &c)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	return c
}

func (c *DragoConfig) WithEnv() *DragoConfig {
	return c
}

func (c *DragoConfig) WithFlags() *DragoConfig {
	return c
}
