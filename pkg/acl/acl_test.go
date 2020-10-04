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

var repo *mockRepository

type mockRepository struct {
	tokens   []*mockToken
	policies []*mockPolicy
}

func (r *mockRepository) FindTokenBySecret(s string) (Token, error) {
	for _, t := range r.tokens {
		if t.secret == s {
			return t, nil
		}
	}
	return nil, nil
}

func (r *mockRepository) GetPolicyByName(n string) (Policy, error) {
	for _, p := range r.policies {
		if p.name == n {
			return p, nil
		}
	}
	return nil, fmt.Errorf("not found : %s", n)
}

type mockToken struct {
	management bool
	secret     string
	policies   []string
}

func (t *mockToken) Policies() []string {
	return t.policies
}

func (t *mockToken) IsManagement() bool {
	return t.management
}

type mockPolicy struct {
	name  string
	rules []*mockRule
}

func (p *mockPolicy) Name() string {
	return p.name
}

func (p *mockPolicy) Rules() []Rule {
	rules := []Rule{}
	for _, r := range p.rules {
		rules = append(rules, r)
	}
	return rules
}

type mockRule struct {
	resource     string
	pattern      string
	capabilities []string
}

func (r *mockRule) Resource() string {
	return r.resource
}

func (r *mockRule) Pattern() string {
	return r.pattern
}

func (r *mockRule) Capabilities() []string {
	return r.capabilities
}

func setupACL() {

	NewResource("namespace").
		AddCapability(capNamespaceReadX, capNamespaceReadY, capNamespaceWriteX, capNamespaceWriteY).
		AddAlias("read", capNamespaceReadX, capNamespaceReadY).
		AddAlias("write", capNamespaceWriteX, capNamespaceWriteY)

	NewResource("network").
		AddCapability(capNetworkRead, capNetworkWrite, capNetworkList).
		AddAlias("read", capNetworkList, capHostRead).
		AddAlias("write", capNetworkWrite, capNetworkRead, capNetworkWrite)

	NewResource("host").
		AddCapability(capNetworkRead, capNetworkWrite, capNetworkList).
		AddAlias("read", capNetworkList, capHostRead).
		AddAlias("write", capNetworkWrite, capNetworkRead, capNetworkWrite)

	PolicyResolver(func(p string) (Policy, error) {
		return repo.GetPolicyByName(p)
	})

	SecretResolver(func(s string) (Token, error) {
		return repo.FindTokenBySecret(s)
	})

	AnonymousToken(&mockToken{false, "", []string{"anonymous"}})
}

func setupMockRepo() {
	repo = &mockRepository{
		tokens: []*mockToken{
			{true, "39076595-19a6-4582-b0d9-bb4a266fd48a", []string{""}},
			{false, "71036287-81d1-474a-b4d5-25c2ee6f57ae", []string{"policy1"}},
			{false, "54c06ace-7da6-443b-a5a2-05da5294fbd5", []string{"policy2"}},
			{false, "e690413b-827b-400e-bc38-92a4b1580eac", []string{"policy1", "policy2"}},
		},
		policies: []*mockPolicy{
			{"policy1", []*mockRule{
				{"namespace", "*", []string{"read"}},
				{"network", "*", []string{"read"}},
				{"host", "", []string{"read"}},
			}},
			{"policy2", []*mockRule{
				{"namespace", "*", []string{"write"}},
				{"network", "*", []string{"deny"}},
				{"host", "*", []string{"list"}},
			}},
			{"policy3", []*mockRule{
				{"network", "*", []string{"read"}},
				{"host", "*", []string{"write"}},
			}},
			{"anonymous", []*mockRule{
				{"network", "*", []string{"list"}},
				{"host", "*", []string{"list"}},
			}},
		},
	}
}

func TestACL(t *testing.T) {

	setupMockRepo()

	t.Run("MissingConfigurations", func(t *testing.T) {

		secret := "54c06ace-7da6-443b-a5a2-05da5294fbd5"

		if _, err := ResolveSecret(secret); err == nil {
			t.Fatal(err)
		}

		// Configure secret resolver
		SecretResolver(func(s string) (Token, error) {
			return repo.FindTokenBySecret(s)
		})
		if _, err := ResolveSecret(secret); err == nil {
			t.Fatal(err)
		}

		// Configure policy resolver
		PolicyResolver(func(p string) (Policy, error) {
			return repo.GetPolicyByName(p)
		})
		if _, err := ResolveSecret(secret); err != nil {
			t.Fatal(err)
		}
	})

	setupACL()

	t.Run("ResolveSecret", func(t *testing.T) {

		secret := "54c06ace-7da6-443b-a5a2-05da5294fbd5"

		acl, err := ResolveSecret(secret)
		if err != nil {
			t.Fatal(err)
		}

		if err := acl.CheckAuthorized("namespace", "12345", "write-x"); err != nil {
			t.Fatalf("%v", err)
		}

		if err := acl.CheckAuthorized("network", "12345", "write"); err == nil {
			t.Fatalf("%v", err)
		}

		if err := acl.CheckAuthorized("host", "12345", "write"); err == nil {
			t.Fatalf("%v", err)
		}
	})

	t.Run("InvalidSecret", func(t *testing.T) {

		secret := "foobar"

		_, err := ResolveSecret(secret)
		if err == nil {
			t.Fatal(err)
		}
	})

	t.Run("ManagementToken", func(t *testing.T) {

		secret := "39076595-19a6-4582-b0d9-bb4a266fd48a"

		acl, err := ResolveSecret(secret)
		if err != nil {
			t.Fatal(err)
		}

		if err := acl.CheckAuthorized("namespace", "12345", "write-x"); err != nil {
			t.Fatalf("%v", err)
		}

		if err := acl.CheckAuthorized("network", "12345", "write"); err != nil {
			t.Fatalf("%v", err)
		}

		if err := acl.CheckAuthorized("host", "12345", "write"); err != nil {
			t.Fatalf("%v", err)
		}
	})

	t.Run("AnonymousToken", func(t *testing.T) {

		secret := ""

		acl, err := ResolveSecret(secret)
		if err != nil {
			t.Fatal(err)
		}

		if err := acl.CheckAuthorized("namespace", "12345", "write-x"); err == nil {
			t.Fatalf("%v", err)
		}

		if err := acl.CheckAuthorized("network", "12345", "write"); err == nil {
			t.Fatalf("%v", err)
		}

		if err := acl.CheckAuthorized("host", "12345", "list"); err != nil {
			t.Fatalf("%v", err)
		}
	})

	t.Run("StringDump", func(t *testing.T) {

		secret := "54c06ace-7da6-443b-a5a2-05da5294fbd5"

		acl, err := ResolveSecret(secret)
		if err != nil {
			t.Fatal(err)
		}

		str := acl.String()
		if str != "" {
			t.Fatal(err)
		}
	})

}
