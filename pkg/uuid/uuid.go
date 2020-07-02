package uuid

import "github.com/google/uuid"

// UUID : Type representing an UUID
type UUID [16]byte

// Nil : Nil UUID
var Nil UUID

// NewRandom : Generate random UUID
func NewRandom() (UUID, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return Nil, err
	}

	return UUID(uuid), nil
}

// String : Convert an UUID to string
func (id UUID) String() string {
	return uuid.UUID(id).String()
}
