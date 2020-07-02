package application

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/seashell/drago/server/domain"
)

// TokenService :
type TokenService interface {
	GetByID(id string) (*domain.Token, error)
	Create(t *domain.Token) (*domain.Token, error)
	Update(t *domain.Token) (*domain.Token, error)
	DeleteByID(id string) (*domain.Token, error)
	FindAll(pageInfo domain.PageInfo) ([]*domain.Token, *domain.Page, error)
	Decode(t *domain.Token) (*domain.Token, error)
}

type tokenService struct {
	hr domain.HostRepository
}

// NewTokenService :
func NewTokenService(hs domain.HostRepository) (TokenService, error) {
	return &tokenService{hs}, nil
}

// GetByID :
func (s *tokenService) GetByID(id string) (*domain.Token, error) {
	return nil, errors.New("Not implemented")
}

// Create :
func (s *tokenService) Create(t *domain.Token) (*domain.Token, error) {

	if t.Type == "client" {
		if t.Subject != "" {

			if _, err := uuid.Parse(t.Subject); err != nil {
				return nil, errors.New("Subject is not a valid UUID")
			}

			exists, err := s.hr.ExistsByID(t.Subject)
			if err != nil {
				return nil, err
			}
			if !exists {
				fmt.Println("Subject does not exist. Issuing token anyways...")
			}
		}
	}

	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	t.ID = uid.String()
	t.IssuedAt = now
	t.ExpiresAt = now + 315360000
	t.NotBefore = now

	secret := os.Getenv("ROOT_SECRET")

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     t.ID,
		"type":   t.Type,
		"sub":    t.Subject,
		"labels": t.Labels,
		"iat":    t.IssuedAt,
		"exp":    t.ExpiresAt,
		"nbf":    t.NotBefore,
	}).SignedString([]byte(secret))

	if err != nil {
		return nil, err
	}

	t.Raw = tokenString

	return t, nil
}

// Update :
func (s *tokenService) Update(t *domain.Token) (*domain.Token, error) {
	return nil, errors.New("Not implemented")
}

// DeleteByID :
func (s *tokenService) DeleteByID(id string) (*domain.Token, error) {
	return nil, errors.New("Not implemented")
}

// FindAll :
func (s *tokenService) FindAll(pageInfo domain.PageInfo) ([]*domain.Token, *domain.Page, error) {
	return nil, nil, errors.New("Not implemented")
}

// Decode :
func (s *tokenService) Decode(t *domain.Token) (*domain.Token, error) {
	secret := os.Getenv("ROOT_SECRET")

	err := t.Decode(secret)

	if err != nil {
		return nil, err
	}
	return t, nil
}
