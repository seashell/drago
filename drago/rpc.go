package drago

import (
	"fmt"

	rpc "github.com/seashell/drago/drago/rpc"
	receiver "github.com/seashell/drago/drago/rpc/receiver"
)

func (s *Server) setupRPCServer() error {

	config := &rpc.ServerConfig{
		Logger:      s.logger,
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
