package concurrent

import (
	"sync"
)

// Map type that can be safely shared between
// goroutines that require read/write access to a map
type Map struct {
	sync.RWMutex
	items map[string]interface{}
}

func NewMap() *Map {
	return &Map{
		items: map[string]interface{}{},
	}
}

// Concurrent map item
type MapItem struct {
	Key   string
	Value interface{}
}

// Sets a key in a concurrent map
func (cm *Map) Set(key string, value interface{}) {
	cm.Lock()
	defer cm.Unlock()
	cm.items[key] = value

}

// Gets a key from a concurrent map
func (cm *Map) Get(key string) (interface{}, bool) {
	cm.RLock()
	defer cm.RUnlock()
	value, ok := cm.items[key]
	return value, ok
}

// Returns the length of a concurrent map
func (cm *Map) Len() int {
	cm.RLock()
	defer cm.RUnlock()
	return len(cm.items)
}

// Deletes a key from a concurrent map
func (cm *Map) Delete(key string) {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.items, key)
}

// Iterates over the items in a concurrent map
// Each item is sent over a channel, so that we can iterate
// over the map using the builtin range keyword
func (cm *Map) Iter() <-chan MapItem {
	c := make(chan MapItem)
	f := func() {
		snapshot := map[string]MapItem{}
		cm.RLock()
		for k, v := range cm.items {
			snapshot[k] = MapItem{Key: k, Value: v}
		}
		cm.RUnlock()
		for _, item := range snapshot {
			c <- item
		}
		close(c)
	}
	go f()

	return c
}
