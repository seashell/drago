package config

import (
	"time"

	"github.com/seashell/drago/pkg/acl"
)

// ACLConfig :
type ACLConfig struct {
	Enabled bool
	// TokenTTL controls for how long we keep ACL tokens in cache.
	TokenTTL time.Duration

	// Model contains the ACL model
	Model *acl.Model
}

// DefaultACLConfig :
func DefaultACLConfig() *ACLConfig {
	return &ACLConfig{
		Enabled:  false,
		TokenTTL: 30 * time.Second,
		Model:    acl.NewModel(),
	}
}
