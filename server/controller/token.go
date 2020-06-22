package controller

import (
	"context"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

// CreateTokenInput :
type CreateTokenInput struct {
	Type    *string `json:"type" validate:"required,oneof='client' 'management'"`
	Subject *string `json:"subject" validate:"required"`
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
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	token := &domain.Token{
		ID:        uid.String(),
		Type:      *in.Type,
		Subject:   *in.Subject,
		IssuedAt:  now,
		ExpiresAt: now + 315360000,
		NotBefore: now,
	}

	secret := os.Getenv("ROOT_SECRET")

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   token.ID,
		"type": token.Type,
		"sub":  *in.Subject,
		"iat":  token.IssuedAt,
		"exp":  token.ExpiresAt,
		"nbf":  token.NotBefore,
	}).SignedString([]byte(secret))

	if err != nil {
		return nil, err
	}

	token.Raw = tokenString

	return token, nil
}

// GetSelfToken :
func (c *Controller) GetSelfToken(ctx context.Context, in *GetSelfTokenInput) (*domain.Token, error) {

	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	claims := jwt.MapClaims{}

	secret := os.Getenv("ROOT_SECRET")

	parser := jwt.Parser{
		SkipClaimsValidation: true,
	}

	_, err = parser.ParseWithClaims(*in.Raw, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	token := &domain.Token{
		ID:        claims["id"].(string),
		Raw:       *in.Raw,
		Type:      claims["type"].(string),
		Subject:   claims["sub"].(string),
		IssuedAt:  int64(claims["iat"].(float64)),
		ExpiresAt: int64(claims["exp"].(float64)),
		NotBefore: int64(claims["nbf"].(float64)),
	}

	return token, nil
}
