package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/seashell/drago/client/nic"
	"github.com/seashell/drago/client/state"
	"github.com/seashell/drago/client/state/boltdb"

	api "github.com/seashell/drago/api"
	structs "github.com/seashell/drago/drago/structs"
	log "github.com/seashell/drago/pkg/log"
)

// Client is the Drago client
type Client struct {
	config *Config
	logger log.Logger

	api *api.Client

	state state.Repository

	niController *nic.Controller

	nodeID         string
	nodeStatus     string
	nodeStatusLock sync.Mutex

	heartbeatTTL  time.Duration
	heartbeatLock sync.Mutex

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

	err := c.init()
	if err != nil {
		c.config.Logger.Errorf("error initializing the client: %v", err)
	}

	// Setup the network controller
	nc, err := nic.NewController(&nic.Config{
		InterfacePrefix: "drago",
		WireguardPath:   c.config.WireguardPath,
	})
	if err != nil {
		return nil, err
	}
	c.niController = nc

	// Register node
	c.registerAndHeartbeat()

	// Start the client!
	go c.run()

	c.logger.Infof("started client (node id %s)", c.NodeID())
	return c, nil
}

func (c *Client) init() error {

	// Ensure the state dir exists
	if c.config.StateDir != "" {
		if err := os.MkdirAll(c.config.StateDir, 0700); err != nil {
			return fmt.Errorf("failed to create state dir: %s", err)
		}
	} else {
		// Otherwise make a temp directory to use.
		p, err := ioutil.TempDir("", "DragoClient")
		if err != nil {
			return fmt.Errorf("failed to create temporary directory: %v", err)
		}

		p, err = filepath.EvalSymlinks(p)
		if err != nil {
			return fmt.Errorf("failed to find temporary directory for the StateDir: %v", err)
		}
		c.config.StateDir = p
	}
	c.logger.Infof("using state directory %s", c.config.StateDir)

	// Create state repository
	repo, err := boltdb.NewStateRepository(path.Join(c.config.StateDir, "client.state"))
	if err != nil {
		return fmt.Errorf("failed to open state database: %v", err)
	}
	c.state = repo

	// Setup the HTTP API client
	cli, err := api.NewClient(&api.Config{
		Address: c.config.Servers[0],
		Token:   c.config.Token,
	})
	if err != nil {
		return fmt.Errorf("failed to setup API client: %v", err)
	}

	c.api = cli

	return nil
}

func (c *Client) watchInterfaces(ch chan *structs.Interface) {

	req := &structs.NodeSpecificRequest{
		ID: c.NodeID(),
	}

	for {
		_, err := c.api.Nodes().GetNodeInterfaces(context.TODO(), req)
		if err != nil {
			// Shutdown often causes EOF errors, so check for shutdown first
			select {
			case <-c.shutdownCh:
				return
			default:
			}

			time.Sleep(c.config.ReconcileInterval)
		}

		// Check for shutdown
		select {
		case <-c.shutdownCh:
			return
		default:
		}

	}
}

// run is a long lived goroutine used to run the client. Shutdown() stops it first
func (c *Client) run() {

	// Watch for changes in interfaces
	interfaceUpdates := make(chan *structs.Interface)
	go c.watchInterfaces(interfaceUpdates)

	for {
		select {
		case update := <-interfaceUpdates:
			// Don't apply updates while shutting down.
			c.shutdownLock.Lock()
			if c.shutdown {
				c.shutdownLock.Unlock()
				return
			}
			// TODO: Apply interface updates
			fmt.Println(update)
			c.shutdownLock.Unlock()

		case <-c.shutdownCh:
			return
		}
	}
}

// registerAndHeartbeat is a long lived goroutine used to register the client
// and then start heartbeating to the server.
func (c *Client) registerAndHeartbeat() {
	// Register the node
	c.retryRegisterNode()

	// Start watching changes for node changes
	// go c.watchNodeUpdates()

	// Start watching for emitting node events
	// go c.watchNodeEvents()

	// Setup the heartbeat timer.
	var heartbeat <-chan time.Time
	heartbeat = time.After(5)

	for {
		select {
		case <-heartbeat:
		case <-c.shutdownCh:
			return
		}
		if err := c.updateNodeStatus(); err != nil {
			if strings.Contains(err.Error(), "node not found") {
				c.logger.Infof("re-registering node")
				c.retryRegisterNode()
				heartbeat = time.After(5 * time.Second)
			} else {
				heartbeat = time.After(5 * time.Second)
			}
		} else {
			c.heartbeatLock.Lock()
			heartbeat = time.After(c.heartbeatTTL)
			c.heartbeatLock.Unlock()
		}
	}
}

func (c *Client) updateNodeStatus() error {

	req := &structs.NodeUpdateStatusRequest{
		ID:     c.NodeID(),
		Status: structs.NodeStatusReady,
	}
	resp, err := c.api.Nodes().UpdateStatus(context.TODO(), req)
	if err != nil {
		return err
	}

	// Update the last heartbeat and the new TTL, capturing the old values
	c.heartbeatLock.Lock()
	c.heartbeatTTL = resp.HeartbeatTTL
	c.heartbeatLock.Unlock()

	return nil
}

func (c *Client) retryRegisterNode() {
	for {
		err := c.registerNode()
		if err == nil {
			return
		}

		c.logger.Errorf("error registering: %v", err)

		select {
		case <-time.After(time.Duration(5 * time.Second)):
		case <-c.shutdownCh:
			return
		}
	}
}

// registerNode is used to register the node or update the registration
func (c *Client) registerNode() error {

	req := &structs.NodeRegisterRequest{
		Node: *c.Node(),
	}

	resp, err := c.api.Nodes().Register(context.TODO(), req)
	if err != nil {
		return err
	}

	c.nodeStatusLock.Lock()
	c.nodeStatus = structs.NodeStatusReady
	c.nodeStatusLock.Unlock()

	c.logger.Infof("node registration complete")

	c.heartbeatLock.Lock()
	defer c.heartbeatLock.Unlock()
	c.heartbeatTTL = resp.HeartbeatTTL

	return nil
}

func (c *Client) Node() *structs.Node {
	return &structs.Node{}
}

// NodeID returns the node ID for the given client
func (c *Client) NodeID() string {
	return c.nodeID
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
