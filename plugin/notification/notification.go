package notification

package main

import (
	"net/rpc"

	log "github.com/seashell/drago/pkg/log"
)

// TODO: Implementation of the drago.Plugin interface, with services exposed through RPC
// Whenever an event occurs in the Drago server (e.g. a new node is registered, or becomes online),
// this plugin interacts with the Drago RPC API in order to obtain information about this event
// and notify users via email, Telegram, slack, etc.

// NotificationPlugin :
type NotificationPlugin struct {
	logger    log.Logger
	rpcServer *rpc.Server
}

// Config : 
type Config struct {
}

// NewNotificationPlugin : Creates a new mesh plugin object parameterized according to the provided configurations.
func NewNotificationPlugin(config *Config) (*NotificationPlugin, error) {
	p := &NotificationPlugin{}
	return p, nil
}

func main() {
	_, err := NewNotificationPlugin(&Config{})
	if err != nil {
		panic(err)
	}
}
