package domain

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

const (
	TokenTypeClient     = "client"
	TokenTypeManagement = "management"
)

type Token struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	Subject   string   `json:"subject"`
	Labels    []string `json:"labels"`
	Raw       string   `json:"secret"`
	IssuedAt  int64    `json:"issuedAt"`
	ExpiresAt int64    `json:"expiresAt"`
	NotBefore int64    `json:"notBefore"`
}

func (t *Token) Decode(secret string) error {

	claims := jwt.MapClaims{}

	parser := jwt.Parser{
		SkipClaimsValidation: true,
	}

	_, err := parser.ParseWithClaims(t.Raw, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}

	t.ID = claims["id"].(string)
	t.Type = claims["type"].(string)
	t.Subject = claims["sub"].(string)
	t.IssuedAt = int64(claims["iat"].(float64))
	t.ExpiresAt = int64(claims["exp"].(float64))
	t.NotBefore = int64(claims["nbf"].(float64))

	iflabels := claims["labels"].([]interface{})
	t.Labels = make([]string, len(iflabels))

	for i, v := range iflabels {
		t.Labels[i] = fmt.Sprint(v)
	}

	return nil
}
