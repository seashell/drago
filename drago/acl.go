package drago

import (
	"context"
	"fmt"

	structs "github.com/seashell/drago/drago/application/structs"
	domain "github.com/seashell/drago/drago/domain"
	"github.com/seashell/drago/pkg/acl"
)

const (
	capBootstrap = "bootstrap"
)

const (
	capTokenList  = "list"
	capTokenRead  = "read"
	capTokenWrite = "write"
)

const (
	capPolicyList  = "list"
	capPolicyRead  = "read"
	capPolicyWrite = "write"
)

func (s *Server) setupACLSystem() error {

	c := &acl.Config{
		Resources: map[string]*acl.Resource{
			"acl": &acl.Resource{
				Name: "acl",
				Capabilities: map[string]*acl.Capability{
					capBootstrap: {capBootstrap, "Allows ACL to be bootstrapped"},
				},
				Aliases: map[string]*acl.CapabilityAlias{},
			},
			"token": &acl.Resource{
				Name: "token",
				Capabilities: map[string]*acl.Capability{
					capTokenRead:  {capTokenRead, "Allows a token to be read"},
					capTokenWrite: {capTokenWrite, "Allows a token to be written"},
					capTokenList:  {capTokenList, "Allows tokens to be listed"},
				},
				Aliases: map[string]*acl.CapabilityAlias{
					"read":  {"read", []string{capTokenRead, capTokenList}},
					"write": {"write", []string{capTokenWrite, capTokenRead, capTokenList}},
				},
			},
			"policy": &acl.Resource{
				Name: "policy",
				Capabilities: map[string]*acl.Capability{
					capPolicyRead:  {capPolicyRead, "Allows a policy to be read"},
					capPolicyWrite: {capPolicyWrite, "Allows a policy to be written"},
					capPolicyList:  {capPolicyList, "Allows policies to be listed"},
				},
				Aliases: map[string]*acl.CapabilityAlias{
					"write": {"read", []string{capPolicyWrite, capPolicyRead, capPolicyList}},
					"read":  {"read", []string{capPolicyRead, capPolicyList}},
				},
			},
		}}

	c.ResolvePolicy = func(n string) (acl.Policy, error) {
		return nil, nil
	}

	c.ResolveToken = func(s string) (acl.Token, error) {
		return nil, nil
	}

	// TODO: replace this with a standard logger
	c.LogActivity = func(s string) {
		fmt.Println("==> LOG: ", s)
	}

	acl.Initialize(c)

	return nil
}

func (s *Server) setupACLPolicies() error {
	defaultPolicies := defaultACLPolicies()
	for _, p := range defaultPolicies {
		in := &structs.ACLPolicyUpsertInput{
			Name:        p.Name,
			Description: p.Description,
		}
		s.services.policies.Upsert(context.Background(), in)
	}
	return nil
}

func defaultACLPolicies() []*domain.ACLPolicy {

	policies := []*domain.ACLPolicy{}

	anonymous := &domain.ACLPolicy{
		Name: "anonymous",
	}

	manager := &domain.ACLPolicy{
		Name: "manager",
	}

	policies = append(policies, anonymous, manager)

	return policies
}
