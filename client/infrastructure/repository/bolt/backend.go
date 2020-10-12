package bolt

import bolt "go.etcd.io/bbolt"

// Backend ...
type Backend struct {
	db *bolt.DB
}

// NewBackend creates a new BoltDB storage backend
func NewBackend(path string) (*Backend, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &Backend{db}, nil

}
