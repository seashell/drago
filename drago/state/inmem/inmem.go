package inmem

import (
	"context"
	"fmt"
	"strings"

	state "github.com/seashell/drago/drago/state"
	concurrent "github.com/seashell/drago/pkg/concurrent"
	log "github.com/seashell/drago/pkg/log"
)

const (
	defaultNamespace = "default"
	defaultPrefix    = "/registry"
)

// StateRepository ...
type StateRepository struct {
	kv     *concurrent.Map
	logger log.Logger
}

// NewStateRepository ...
func NewStateRepository(logger log.Logger) *StateRepository {
	return &StateRepository{
		kv:     concurrent.NewMap(),
		logger: logger,
	}
}

// Name ...
func (b *StateRepository) Name() string {
	return "inmem"
}

type transaction struct {
}

func (t transaction) Commit() (interface{}, error) {
	return nil, nil
}

// Transaction ...
func (b *StateRepository) Transaction(ctx context.Context) state.Transaction {
	return transaction{}
}

// Dump ...
func (b *StateRepository) Dump() <-chan string {
	dumpCh := make(chan string, 1)
	padding := 72
	go func() {
		for el := range b.kv.Iter() {
			dumpCh <- fmt.Sprintf("%s%s %T", el.Key, strings.Repeat(" ", padding-len(el.Key)), el.Value)
		}
		close(dumpCh)
	}()
	return dumpCh
}

// Clear ...
func (b *StateRepository) Clear() {
	b.kv = &concurrent.Map{}
}

func resourcePrefix(resourceType string) string {
	return fmt.Sprintf("%s/%s/%s", defaultPrefix, resourceType, defaultNamespace)
}

func resourceKey(resourceType, resourceID string) string {
	return fmt.Sprintf("%s/%s", resourcePrefix(resourceType), resourceID)
}
