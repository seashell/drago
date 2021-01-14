package boltdb

import (
	"encoding/json"

	"github.com/seashell/drago/drago/structs"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
)

var (
	interfacesBucketName    = []byte("interfaces")
	interfaceKeysBucketName = []byte("keys")
)

// StateRepository ...
type StateRepository struct {
	db *bolt.DB
}

// NewStateRepository creates a new BoltDB state repository
func NewStateRepository(path string) (*StateRepository, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists(interfacesBucketName)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(interfaceKeysBucketName)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return &StateRepository{db}, nil

}

func (r *StateRepository) Name() string {
	return "boltdb"
}

func (r *StateRepository) Interfaces() ([]*structs.Interface, error) {

	ifaces := []*structs.Interface{}

	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(interfacesBucketName)
		err := b.ForEach(func(k []byte, v []byte) error {

			iface := &structs.Interface{}

			err := decode(v, iface)
			if err != nil {
				return err
			}

			ifaces = append(ifaces, iface)

			return nil
		})
		return err
	})

	return ifaces, err

}

func (r *StateRepository) UpsertInterface(iface *structs.Interface) error {

	err := r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(interfacesBucketName)
		return b.Put([]byte(iface.ID), encode(iface))
	})

	return err
}

func (r *StateRepository) DeleteInterfaces(ids []string) error {

	err := r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(interfacesBucketName)
		for _, id := range ids {
			if err := b.Delete([]byte(id)); err != nil {
				return err
			}
		}
		b = tx.Bucket(interfaceKeysBucketName)
		for _, id := range ids {
			if err := b.Delete([]byte(id)); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (r *StateRepository) InterfaceKeyByID(id string) (string, error) {

	key := ""

	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(interfaceKeysBucketName)

		v := b.Get([]byte(id))

		err := decode(v, &key)
		if err != nil {
			return err
		}

		return nil
	})

	return key, err

}

func (r *StateRepository) UpsertInterfaceKey(id, key string) error {

	err := r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(interfaceKeysBucketName)
		return b.Put([]byte(id), encode(key))
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
