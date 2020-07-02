package controller

import (
	"context"

	"github.com/pkg/errors"
	"github.com/seashell/drago/server/controller/pagination"
	"github.com/seashell/drago/server/domain"
)

// GetTokenInput :
type GetTokenInput struct {
	ID *string `json:"type" validate:"required,uuid4"`
}

// CreateTokenInput :
type CreateTokenInput struct {
	Type    *string  `json:"type" validate:"required,oneof='client' 'management'"`
	Subject *string  `json:"subject" validate:"required"`
	Labels  []string `json:"labels" validate:"dive,omitempty,min=1,max=64"`
}

// DeleteTokenInput :
type DeleteTokenInput struct {
	ID string `validate:"required,uuid4"`
}

// ListTokensInput :
type ListTokensInput struct {
	pagination.Input
	NetworkIDFilter string `query:"networkId" validate:"omitempty,uuid4"`
}

// GetSelfTokenInput :
type GetSelfTokenInput struct {
	Raw *string `json:"type" validate:"required"`
}

// CreateToken :
func (c *Controller) CreateToken(ctx context.Context, in *CreateTokenInput) (*domain.Token, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	if in.Labels == nil {
		in.Labels = []string{}
	}

	t := &domain.Token{
		Type:    *in.Type,
		Subject: *in.Subject,
		Labels:  in.Labels,
	}

	res, err := c.ts.Create(t)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// GetSelfToken :
func (c *Controller) GetSelfToken(ctx context.Context, in *GetSelfTokenInput) (*domain.Token, error) {

	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	t := &domain.Token{
		Raw: *in.Raw,
	}

	res, err := c.ts.Decode(t)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}
