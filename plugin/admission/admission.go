package main

import (
	"log"
	"net/rpc"
)

type AdmissionPlugin struct {
	logger    log.Logger
	rpcServer *rpc.Server
}

type Config struct {
}

// Creates a new admission plugin object parameterized according to the provided configurations.
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
