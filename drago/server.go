package drago

import (
	"context"
	"sync"

	application "github.com/seashell/drago/drago/application"
	rpc "github.com/seashell/drago/drago/infrastructure/rpc"
	receiver "github.com/seashell/drago/drago/infrastructure/rpc/receiver"
	boltdb "github.com/seashell/drago/drago/infrastructure/storage/boltdb"
	log "github.com/seashell/drago/pkg/log"
	logrus "github.com/seashell/drago/pkg/log/logrus"
)

type Server struct {
	config *Config
	logger log.Logger

	shutdown     bool
	shutdownLock sync.Mutex

	rpcServer *rpc.Server

	services struct {
		networks application.NetworkService
	}

	shutdownCtx    context.Context
	shutdownCancel context.CancelFunc
	shutdownCh     <-chan struct{}
}

// Create a new Drago server, potentially returning an error
func NewServer(config *Config) (*Server, error) {

	logger, err := logrus.NewLoggerAdapter(logrus.Config{
		LoggerOptions: log.LoggerOptions{
			Prefix: "drago: ",
			Level:  logrus.Debug,
		},
	})
	if err != nil {
		panic(err)
	}

	s := &Server{
		config: config,
		logger: logger,
	}

	s.shutdownCtx, s.shutdownCancel = context.WithCancel(context.Background())
	s.shutdownCh = s.shutdownCtx.Done()

	backend, err := boltdb.NewBackend("database.db")
	if err != nil {
		panic(err)
	}

	networkRepository := boltdb.NewNetworkRepositoryAdapter(backend)

	s.services.networks = application.NewNetworkService(networkRepository)

	// Setup RPC server
	if err := s.setupRPCServer(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) setupRPCServer() error {

	config := &rpc.ServerConfig{
		Logger:      s.logger,
		BindAddress: s.config.BindAddr,
		Receivers: map[string]interface{}{
			"Networks": receiver.NewNetworkReceiverAdapter(s.services.networks),
		},
	}

	rpcServer, err := rpc.NewServer(config)
	if err != nil {
		return err
	}

	s.rpcServer = rpcServer

	s.rpcServer.Run()

	return nil
}

// Run
func (s *Server) Run() {}
