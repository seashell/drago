package controller

import (
	"context"
	"time"

	"gopkg.in/jeevatkm/go-model.v1"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

// SynchronizeHostInput :
type SynchronizeHostInput struct {
	ID string `validate:"required,uuid4"`
}

// GetHostSettingsInput :
type GetHostSettingsInput struct {
	ID string `validate:"required,uuid4"`
}

// UpdateHostStateInput :
type UpdateHostStateInput struct {
	ID         string `validate:"required,uuid4"`
	Interfaces []*struct {
		Name      *string `json:"name" validate:"required"`
		PublicKey *string `json:"publicKey" validate:"required"`
	} `json:"interfaces"`
	Peers []*struct {
		LatencyMs     uint64    `json:"latencyMs,omitempty"`
		LastHandshake time.Time `json:"lastHandshake,omitempty"`
	} `json:"peers"`
}

// GetHostSettingsByID :
func (c *Controller) GetHostSettings(ctx context.Context, in *GetHostSettingsInput) (*domain.HostSettings, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}
	settings, err := c.ss.GetHostSettingsByID(in.ID)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}
	return settings, nil
}

// UpdateHostState :
func (c *Controller) UpdateHostState(ctx context.Context, in *UpdateHostStateInput) (*domain.HostState, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	stateIn := &domain.HostState{
		Interfaces: []*domain.WgInterfaceState{},
		Peers:      []*domain.WgPeerState{},
	}

	for _, iface := range in.Interfaces {
		stateIn.Interfaces = append(stateIn.Interfaces, &domain.WgInterfaceState{
			Name:      iface.Name,
			PublicKey: iface.PublicKey,
		})
	}

	// for _, peers := range in.Peers {
	//	stateIn.Peers = append(stateIn.Peers, &domain.WgPeerState{})
	// }

	stateOut, err := c.ss.UpdateHostState(in.ID, stateIn)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}
	return stateOut, nil
}

// SynchronizeHost :
func (c *Controller) SynchronizeHost(ctx context.Context, in *SynchronizeHostInput) (*domain.HostSettings, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	state := &domain.HostState{}

	errs := model.Copy(state, in)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	settings, err := c.ss.SynchronizeHost(in.ID, state)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}
	return settings, nil
}
