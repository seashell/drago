package boltdb

import (
	"encoding/json"

	"github.com/seashell/drago/drago/structs"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
)

var (
	interfacesBucketName = []byte("interfaces")
	peersBucketName      = []byte("peers")
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

		_, err = tx.CreateBucketIfNotExists(peersBucketName)
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
		return b.Put(interfacesBucketName, encode(iface))
	})

	return err
}

func (r *StateRepository) Peers() ([]*structs.Peer, error) {
	peers := []*structs.Peer{}

	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(peersBucketName)
		err := b.ForEach(func(k []byte, v []byte) error {

			peer := &structs.Peer{}

			err := decode(v, peer)
			if err != nil {
				return err
			}

			peers = append(peers, peer)

			return nil
		})
		return err
	})

	return peers, err
}

func (r *StateRepository) UpsertPeer(peer *structs.Peer) error {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(peersBucketName)
		return b.Put(peersBucketName, encode(peer))
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
