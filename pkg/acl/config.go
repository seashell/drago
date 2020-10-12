package acl

import (
	"github.com/seashell/drago/pkg/log"
)

// ResolverConfig contains configurations for an ACL Resolver.
type ResolverConfig struct {
	Logger         log.Logger
	Model          *Model
	SecretResolver SecretResolverFunc
	PolicyResolver PolicyResolverFunc
}

// DefaultResolverConfig returns default configurations
// for an ACL Resolver.
func DefaultResolverConfig() *ResolverConfig {
	return &ResolverConfig{
		Logger: &simpleLogger{
			fields: map[string]interface{}{},
			options: log.LoggerOptions{
				Prefix: "acl",
				Level:  levelInfo,
			},
		},
	}
}

// Merge merges two ResolverConfig structs, returning the result.
func (c *ResolverConfig) Merge(in *ResolverConfig) *ResolverConfig {
	res := *c
	if in.Logger != nil {
		res.Logger = in.Logger
	}
	if in.Model != nil {
		res.Model = in.Model
	}
	if in.SecretResolver != nil {
		res.SecretResolver = in.SecretResolver
	}
	if in.PolicyResolver != nil {
		res.PolicyResolver = in.PolicyResolver
	}
	return &res
}
