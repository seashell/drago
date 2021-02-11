package main

import (
	"log"
	"net/rpc"
)

// AdmissionPlugin :
type AdmissionPlugin struct {
	logger    log.Logger
	rpcServer *rpc.Server
}

// Config :
type Config struct {
}

// NewAdmissionPlugin : Creates a new admission plugin object parameterized according to the provided configurations.
func NewAdmissionPlugin(config *Config) (*AdmissionPlugin, error) {
	p := &AdmissionPlugin{}
	return p, nil
}

func main() {
	_, err := NewAdmissionPlugin(&Config{})
	if err != nil {
		panic(err)
	}
}
