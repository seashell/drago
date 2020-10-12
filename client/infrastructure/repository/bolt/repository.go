package bolt

import (
	"encoding/json"
	"fmt"

	"github.com/seashell/drago/client/domain"
	"go.etcd.io/bbolt"
)

const (
	stateBucket = "state"
	stateKey    = "state"
)

// RepositoryAdapter :
type RepositoryAdapter struct {
	backend *Backend
}

// NewStateRepositoryAdapter :
func NewStateRepositoryAdapter(backend *Backend) domain.Repository {

	err := backend.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(stateBucket))
		return err
	})

	if err != nil {
		panic(err)
	}

	return &RepositoryAdapter{
		backend: backend,
	}
}

// Get :
func (a *RepositoryAdapter) Get() (*domain.Host, error) {

	state := &domain.Host{}

	err := a.backend.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(stateBucket))
		v := b.Get([]byte(stateKey))
		if v == nil {
			return fmt.Errorf("not found")
		}
		return decode(v, state)
	})

	if err != nil {
		return nil, err
	}

	return state, nil
}

// Save :
func (a *RepositoryAdapter) Save(state *domain.Host) error {

	err := a.backend.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(stateBucket))
		return b.Put([]byte(stateKey), encode(state))
	})

	return err
}

func encode(in interface{}) []byte {
	out, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return out
}

func decode(encoded []byte, out interface{}) error {
	return json.Unmarshal(encoded, out)
}
