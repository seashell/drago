package acl

import (
	"fmt"
	"testing"
)

const (
	capNamespaceReadX  = "read-x"
	capNamespaceReadY  = "read-y"
	capNamespaceWriteX = "write-x"
	capNamespaceWriteY = "write-y"
)

const (
	capNetworkRead  = "read"
	capNetworkWrite = "write"
	capNetworkList  = "list"
)

const (
	capHostRead  = "read"
	capHostWrite = "write"
	capHostList  = "list"
)

var c = &Config{
	Resources: map[string]*Resource{
		"namespace": &Resource{
			Name: "namespace",
			Capabilities: map[string]*Capability{
				capNamespaceReadX:  {capNamespaceReadX, "Allows x to be read on a namespace"},
				capNamespaceReadY:  {capNamespaceReadY, "Allows y to be read on a namespace"},
				capNamespaceWriteX: {capNamespaceWriteX, "Allows x to be written on a namespace"},
				capNamespaceWriteY: {capNamespaceWriteY, "Allows y to be written on a namespace"},
			},
			Aliases: map[string]*CapabilityAlias{
				"read":  {"read", []string{capNamespaceReadX, capNamespaceReadY}},
				"write": {"write", []string{capNamespaceWriteX, capNamespaceWriteY}},
			},
		},
		"network": &Resource{
			Name: "network",
			Capabilities: map[string]*Capability{
				capNetworkRead:  {capNetworkRead, "Allows a network to be read"},
				capNetworkWrite: {capNetworkWrite, "Allows a network to be written"},
				capNetworkList:  {capNetworkList, "Allows networks to be listed"},
			},
			Aliases: map[string]*CapabilityAlias{
				"write": {"read", []string{capNetworkWrite, capNetworkRead, capNetworkList}},
				"read":  {"read", []string{capNetworkRead, capNetworkList}},
			},
		},
		"host": &Resource{
			Name: "host",
			Capabilities: map[string]*Capability{
				capHostRead:  {capHostRead, "Allows a host to be read"},
				capHostWrite: {capHostWrite, "Allows a host to be written"},
				capHostList:  {capHostList, "Allows hosts to be listed"},
			},
			Aliases: map[string]*CapabilityAlias{
				"write": {"read", []string{capHostWrite, capHostRead, capHostList}},
				"read":  {"read", []string{capHostRead, capHostList}},
			},
		},
	}}

var mockRepo = &repository{
	tokens: []*token{
		{"71036287-81d1-474a-b4d5-25c2ee6f57ae", []string{"policy1"}},
		{"54c06ace-7da6-443b-a5a2-05da5294fbd5", []string{"policy2"}},
		{"e690413b-827b-400e-bc38-92a4b1580eac", []string{"policy1", "policy2"}},
	},
	policies: []*policy{
		{"policy1", []Rule{
			&rule{"namespace", "*", []string{"read"}},
			&rule{"network", "*", []string{"read"}},
			&rule{"host", "", []string{"read"}},
		}},
		{"policy2", []Rule{
			&rule{"namespace", "*", []string{"write"}},
			&rule{"network", "*", []string{"write"}},
			&rule{"host", "*", []string{"list"}},
		}},
		{"policy3", []Rule{
			&rule{"network", "*", []string{"read"}},
			&rule{"host", "*", []string{"write"}},
		}},
	},
}

func TestACL(t *testing.T) {

	c.ResolvePolicy = func(n string) (Policy, error) {
		return mockRepo.GetPolicyByName(n)
	}

	c.ResolveToken = func(s string) (Token, error) {
		return mockRepo.FindTokenBySecret(s)
	}

	// TODO: replace this with a standard logger
	c.LogActivity = func(s string) {
		fmt.Println("==> LOG: ", s)
	}

	// Initialize ACL
	Initialize(c)

	// Build ACL based on a secret
	acl, err := NewACLFromSecret("54c06ace-7da6-443b-a5a2-05da5294fbd5")
	if err != nil {
		t.Fatal(err)
	}

	// Print a string representation of the ACL.
	fmt.Println(acl.String())

	// Check if the ACL is authorized to perform a "write" operation
	// on the resource of type "network" and identified by "abc".
	fmt.Printf("isAuthorized(write, network, abcxy) = %v\n", acl.IsAuthorized("network", "abc", "write"))

	// Check if the ACL is authorized to perform a "write" operation
	// on the resource of type "host" and identified by "123".
	fmt.Printf("isAuthorized(write, host, abcxy) = %v\n", acl.IsAuthorized("host", "123", "write"))
}
