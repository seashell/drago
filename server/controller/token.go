package controller

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// CreateTokenInput :
type CreateTokenInput struct {
	Type    *string `json:"type" validate:"required,oneof='client' 'management'"`
	Subject *string `json:"subject" validate:"required"`
}

// CreateToken :
func (c *Controller) CreateToken(ctx context.Context, in *CreateTokenInput) (*jwt.Token, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	claims := jwt.StandardClaims{
		Subject: *in.Subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tkn, err := token.SignedString("SEASHELLSECRET")
	if err != nil {
		return nil, err
	}
	token.Raw = tkn

	return token, nil
}
