package inmem

import (
	"fmt"
	"strings"
)

const (
	defaultNamespace = "default"
	defaultPrefix    = "/registry"
)

// Backend ...
type Backend struct {
	kv map[string]interface{}
}

// NewBackend ...
func NewBackend() *Backend {
	return &Backend{
		kv: map[string]interface{}{},
	}
}

// Dump ...
func (b *Backend) Dump() <-chan string {
	dumpCh := make(chan string, 1)
	padding := 72
	go func() {
		for k, v := range b.kv {
			dumpCh <- fmt.Sprintf("%s%s %T", k, strings.Repeat(" ", padding-len(k)), v)
		}
		close(dumpCh)
	}()
	return dumpCh
}

// Clear ...
func (b *Backend) Clear() {
	b.kv = map[string]interface{}{}
}

func resourcePrefix(resourceType string) string {
	return fmt.Sprintf("%s/%s/%s", defaultPrefix, resourceType, defaultNamespace)
}

func resourceKey(resourceType, resourceID string) string {
	return fmt.Sprintf("%s/%s/%s/%s", defaultPrefix, resourceType, defaultNamespace, resourceID)
}
