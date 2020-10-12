package client

import (
	"context"
	"math/rand"
	"path"
	"sync"
	"time"

	application "github.com/seashell/drago/client/application"
	drago "github.com/seashell/drago/client/infrastructure/gateway/drago"
	nic "github.com/seashell/drago/client/infrastructure/gateway/nic"
	bolt "github.com/seashell/drago/client/infrastructure/repository/bolt"
	log "github.com/seashell/drago/pkg/log"
)

// Client is the Drago client
type Client struct {
	config   *Config
	logger   log.Logger
	services struct {
		reconciliation application.ReconciliationService
	}
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// New is used to create a new Drago client from the
// configuration, potentially returning an error
func New(config *Config) (*Client, error) {

	rand.Seed(time.Now().Unix())

	config = DefaultConfig().Merge(config)

	c := &Client{
		config: config,
		logger: config.Logger,
	}

	err := c.setupApplication()
	if err != nil {
		return nil, err
	}

	go func() {
		ctx := context.TODO()
		for {
			c.services.reconciliation.Reconcile(ctx)
			time.Sleep(c.config.ReconcileInterval)
		}
	}()

	return c, nil
}

func (c *Client) setupApplication() error {

	gw, err := drago.NewDragoGatewayAdapter(c.config.Servers[0], c.config.Token)
	if err != nil {
		return err
	}

	nc, err := nic.NewController(&nic.Config{
		InterfacePrefix: "drago",
		WireguardPath:   c.config.WireguardPath,
	})
	if err != nil {
		return err
	}

	backend, err := bolt.NewBackend(path.Join(c.config.StateDir, "client.state"))
	if err != nil {
		return err
	}

	config := &application.Config{
		DragoGateway:        gw,
		InterfaceController: nc,
		StateRepository:     bolt.NewStateRepositoryAdapter(backend),
	}

	c.services.reconciliation = application.NewReconciliationService(config)

	return nil
}

// Shutdown is used to tear down the client
func (c *Client) Shutdown() error {
	c.shutdownLock.Lock()
	defer c.shutdownLock.Unlock()

	if c.shutdown {
		c.logger.Infof("client already shutdown")
		return nil
	}
	c.logger.Infof("shutting down")

	c.shutdown = true
	close(c.shutdownCh)

	return nil
}
