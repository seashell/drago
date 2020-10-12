package drago

import (
	"fmt"

	rpc "github.com/seashell/drago/drago/infrastructure/delivery/rpc"
	receiver "github.com/seashell/drago/drago/infrastructure/delivery/rpc/receiver"
)

func (s *Server) setupRPCServer() error {

	config := &rpc.ServerConfig{
		Logger:      s.config.Logger,
		BindAddress: fmt.Sprintf("%s:%d", s.config.BindAddr, s.config.Ports.RPC),
		Receivers: map[string]interface{}{
			"Status": receiver.NewStatusReceiverAdapter(),
			// "ACL": receiver.NewACLReceiverAdapter(),
			// "Hosts": receiver.NewHostReceiverAdapter(),
			// "Networks": receiver.NewNetworkReceiverAdapter(),
		},
	}

	rpcServer, err := rpc.NewServer(config)
	if err != nil {
		return err
	}

	s.rpcServer = rpcServer

	return nil
}
