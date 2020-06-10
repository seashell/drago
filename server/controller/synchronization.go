package controller

import (
	"context"

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
	ID string `validate:"required,uuid4"`
}

// GetHostSettingsByID :
func (c *Controller) GetHostSettingsByID(ctx context.Context, in *GetHostSettingsInput) (*domain.HostSettings, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}
	settings, err := c.ss.GetHostSettingsByID(in.ID)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}
	return settings, ErrNotImplemented
}

// UpdateHostStateByID :
func (c *Controller) UpdateHostStateByID(ctx context.Context, in *UpdateHostStateInput) error {
	err := c.v.Struct(in)
	if err != nil {
		return errors.Wrap(ErrInvalidInput, err.Error())
	}
	return ErrNotImplemented
}

// SynchronizeHost :
func (c *Controller) SynchronizeHost(ctx context.Context, in *SynchronizeHostInput) (*domain.HostSettings, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}
	return nil, ErrNotImplemented
}
