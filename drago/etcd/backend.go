package etcd

import (
	"go.etcd.io/etcd/clientv3"
)

const (
	defaultNamespace = "default"
	defaultPrefix    = "/registry"
)

// Backend :
type Backend struct {
	client *clientv3.Client
}

// NewBackend :
func NewBackend(client *clientv3.Client) (*Backend, error) {
	backend := &Backend{client: client}
	return backend, nil
}
