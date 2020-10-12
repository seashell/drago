package receiver

import (
	"github.com/seashell/drago/drago/application/structs"
)

// StatusReceiverAdapter is used to check on server status
type StatusReceiverAdapter struct{}

// NewStatusReceiverAdapter :
func NewStatusReceiverAdapter() *StatusReceiverAdapter {
	return &StatusReceiverAdapter{}
}

// Ping is used to check for connectivity
func (r *StatusReceiverAdapter) Ping(args struct{}, reply *struct{}) error {
	return nil
}

// Version returns the version of the server
func (r *StatusReceiverAdapter) Version(in struct{}, out *structs.StatusVersionOutput) error {
	return nil
}
