package main

type LeasePlugin struct {
	logger    log.Logger
	rpcServer *rpc.Server
}

type Config struct {
}

// Creates a new lease plugin object parameterized according to the provided configurations.
func NewLeasePlugin(config *Config) (*MeshPlugin, error) {
	p := &MeshPlugin{}
	return p, nil
}

func main() {
	_, err := NewLeasePlugin(&Config{})
	if err != nil {
		panic(err)
	}
}
