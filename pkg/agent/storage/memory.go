package storage

import (
	"errors"
	"fmt"
	"time"
)

type inMemoryStore struct {
	nodes map[int]*Node
	edges map[int]*Edge
	seq   int
	seqe  int
}

func NewInMemoryStore() (*inMemoryStore, error) {
	return &inMemoryStore{
		nodes: make(map[int]*Node, 0),
		edges: make(map[int]*Edge, 0),
		seq:   1,
		seqe:  1,
	}, nil
}

// ********************************
// ***********  NODES  ************
// ********************************

func (s *inMemoryStore) SelectAllNodes() (map[int]*Node, error) {
	return s.nodes, nil
}

func (s *inMemoryStore) SelectNode(id int) (*Node, error) {
	if n, ok := s.nodes[id]; ok {
		return n, nil
	}
	return nil, errors.New("Object not found")
}

func (s *inMemoryStore) DeleteNode(id int) error {
	if _, found := s.nodes[id]; found {
		delete(s.nodes, id)
		return nil
	}
	return errors.New("Object not found")
}

func (s *inMemoryStore) UpdateNode(id int, u *Node) (*Node, error) {

	n, found := s.nodes[id]

	if !found {
		return nil, errors.New("Object not found")
	}

	n.PublicKey = u.PublicKey
	n.UpdatedAt = time.Now()

	return n, nil
}

func (s *inMemoryStore) InsertNode(n *Node) (*Node, error) {

	for _, v := range s.nodes {
		if v.Name == n.Name {
			return nil, errors.New("Attribute 'name' must be unique")
		}
	}

	n.ID = s.seq
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	s.nodes[s.seq] = n

	s.seq += 1

	return n, nil
}

// ********************************
// ***********  PEERS  ************
// ********************************

// TODO: move part of the logic below to the API level
func (s *inMemoryStore) SelectAllPeersForNode(id int) ([]*Peer, error) {

	peers := make([]*Peer, 0)
	fmt.Println("id", id)
	for _, v := range s.edges {
		if v.Source == id {
			peers = append(peers, &Peer{
				Endpoint:            s.nodes[v.Target].AdvertiseAddr + s.nodes[v.Target].ListenPort,
				AllowedIPs:          s.nodes[v.Target].Address,
				PublicKey:           s.nodes[v.Target].PublicKey,
				PersistentKeepalive: v.PersistentKeepalive,
			})
		} else if v.Target == id {
			peers = append(peers, &Peer{
				Endpoint:            s.nodes[v.Source].AdvertiseAddr + s.nodes[v.Source].ListenPort,
				AllowedIPs:          s.nodes[v.Source].Address,
				PublicKey:           s.nodes[v.Source].PublicKey,
				PersistentKeepalive: v.PersistentKeepalive,
			})
		}

	}
	return peers, nil
}

// ********************************
// ***********  EDGES  ************
// ********************************

func (s *inMemoryStore) SelectAllEdges() (map[int]*Edge, error) {
	return s.edges, nil
}

func (s *inMemoryStore) SelectEdge(id int) (*Edge, error) {
	if e, ok := s.edges[id]; ok {
		return e, nil
	}
	return nil, errors.New("Object not found")
}

func (s *inMemoryStore) DeleteEdge(id int) error {
	if _, found := s.edges[id]; found {
		delete(s.edges, id)
		return nil
	}
	return errors.New("Object not found")
}

func (s *inMemoryStore) UpdateEdge(id int, u *Edge) (*Edge, error) {

	e, found := s.edges[id]

	if !found {
		return nil, errors.New("Object not found")
	}

	return e, nil
}

func (s *inMemoryStore) InsertEdge(e *Edge) (*Edge, error) {

	for _, v := range s.edges {
		if v.Source == e.Source && v.Target == e.Target {
			return nil, errors.New("Redundant edge")
		} else if v.Source == e.Target && v.Target == e.Source {
			return nil, errors.New("Redundant edge")
		}
	}

	e.ID = s.seqe
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()

	s.edges[s.seqe] = e

	s.seqe += 1

	return e, nil
}
