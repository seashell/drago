package inmem

import (
	"fmt"
	"strings"
	"sync"

	log "github.com/seashell/drago/pkg/log"
)

const (
	defaultNamespace = "default"
	defaultPrefix    = "/registry"
)

// StateRepository ...
type StateRepository struct {
	kv     map[string]interface{}
	logger log.Logger
	mu     sync.RWMutex
}

// NewStateRepository ...
func NewStateRepository(logger log.Logger) *StateRepository {
	return &StateRepository{
		kv:     map[string]interface{}{},
		logger: logger,
	}
}

// Name ...
func (b *StateRepository) Name() string {
	return "inmem"
}

// Dump ...
func (b *StateRepository) Dump() <-chan string {
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
func (b *StateRepository) Clear() {
	b.kv = map[string]interface{}{}
}

func resourcePrefix(resourceType string) string {
	return fmt.Sprintf("%s/%s/%s", defaultPrefix, resourceType, defaultNamespace)
}

func resourceKey(resourceType, resourceID string) string {
	return fmt.Sprintf("%s/%s", resourcePrefix(resourceType), resourceID)
}
