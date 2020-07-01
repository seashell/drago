package uuid

import "github.com/google/uuid"

type UUID [16]byte

var Nil UUID

func NewRandom() (UUID, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return Nil, err
	}

	return UUID(uuid), nil
}

func (id UUID) String() string {
	return uuid.UUID(id).String()
}
