package storage

import (
	"errors"
	"fmt"
	"time"
)

type inMemoryStore struct {
	hosts map[int]*Host
	links map[int]*Link
	seq   int
	seql  int
}

func NewInMemoryStore() (*inMemoryStore, error) {
	return &inMemoryStore{
		hosts: make(map[int]*Host, 0),
		links: make(map[int]*Link, 0),
		seq:   1,
		seql:  1,
	}, nil
}

// ********************************
// ***********  HOSTS  ************
// ********************************

func (s *inMemoryStore) SelectAllHosts() (map[int]*Host, error) {
	return s.hosts, nil
}

func (s *inMemoryStore) SelectHost(id int) (*Host, error) {
	if h, ok := s.hosts[id]; ok {
		return h, nil
	}
	return nil, errors.New("Object not found")
}

func (s *inMemoryStore) DeleteHost(id int) error {
	if _, found := s.hosts[id]; found {
		delete(s.hosts, id)
		return nil
	}
	return errors.New("Object not found")
}

func (s *inMemoryStore) UpdateHost(id int, uh *Host) (*Host, error) {

	h, found := s.hosts[id]

	if !found {
		return nil, errors.New("Object not found")
	}

	h.PublicKey = uh.PublicKey
	h.UpdatedAt = time.Now()

	return h, nil
}

func (s *inMemoryStore) InsertHost(h *Host) (*Host, error) {

	for _, v := range s.hosts {
		if v.Name == h.Name {
			return nil, errors.New("Attribute 'name' must be unique")
		}
	}

	h.ID = s.seq
	h.CreatedAt = time.Now()
	h.UpdatedAt = time.Now()

	s.hosts[s.seq] = h

	s.seq += 1

	return h, nil
}

// ********************************
// ***********  PEERS  ************
// ********************************

// TODO: move part of the logic below to the API level
func (s *inMemoryStore) SelectAllPeersForHost(id int) ([]*Peer, error) {

	peers := make([]*Peer, 0)
	fmt.Println("id", id)
	for _, v := range s.links {
		if v.Source == id {
			peers = append(peers, &Peer{
				Endpoint:            s.hosts[v.Target].AdvertiseAddr + s.hosts[v.Target].ListenPort,
				AllowedIPs:          s.hosts[v.Target].Address,
				PublicKey:           s.hosts[v.Target].PublicKey,
				PersistentKeepalive: v.PersistentKeepalive,
			})
		} else if v.Target == id {
			peers = append(peers, &Peer{
				Endpoint:            s.hosts[v.Source].AdvertiseAddr + s.hosts[v.Source].ListenPort,
				AllowedIPs:          s.hosts[v.Source].Address,
				PublicKey:           s.hosts[v.Source].PublicKey,
				PersistentKeepalive: v.PersistentKeepalive,
			})
		}

	}
	return peers, nil
}

// ********************************
// ***********  LINKS  ************
// ********************************

func (s *inMemoryStore) SelectAllLinks() (map[int]*Link, error) {
	return s.links, nil
}

func (s *inMemoryStore) SelectLink(id int) (*Link, error) {
	if l, ok := s.links[id]; ok {
		return l, nil
	}
	return nil, errors.New("Object not found")
}

func (s *inMemoryStore) DeleteLink(id int) error {
	if _, found := s.links[id]; found {
		delete(s.links, id)
		return nil
	}
	return errors.New("Object not found")
}

func (s *inMemoryStore) UpdateLink(id int, ul *Link) (*Link, error) {

	l, found := s.links[id]

	if !found {
		return nil, errors.New("Object not found")
	}

	return l, nil
}

func (s *inMemoryStore) InsertLink(l *Link) (*Link, error) {

	for _, v := range s.links {
		if v.Source == l.Source && v.Target == l.Target {
			return nil, errors.New("Redundant link")
		} else if v.Source == l.Target && v.Target == l.Source {
			return nil, errors.New("Redundant link")
		}
	}

	l.ID = s.seql
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()

	s.links[s.seql] = l

	s.seql += 1

	return l, nil
}
