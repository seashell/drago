package acl

import (
	"context"
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
var ctx = context.TODO()

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
	privileged bool
	secret     string
	policies   []string
}

func (t *mockToken) Policies() []string {
	return t.policies
}

func (t *mockToken) IsPrivileged() bool {
	return t.privileged
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

func ACLResolverConfig() *ResolverConfig {

	model := NewModel()
	model.Resource("namespace").
		Capabilities(capNamespaceReadX, capNamespaceReadY, capNamespaceWriteX, capNamespaceWriteY).
		Alias("read", capNamespaceReadX, capNamespaceReadY).
		Alias("write", capNamespaceWriteX, capNamespaceWriteY)
	model.Resource("network").
		Capabilities(capNetworkRead, capNetworkWrite, capNetworkList).
		Alias("read", capNetworkList, capHostRead).
		Alias("write", capNetworkWrite, capNetworkRead, capNetworkWrite)
	model.Resource("host").
		Capabilities(capNetworkRead, capNetworkWrite, capNetworkList).
		Alias("read", capNetworkList, capHostRead).
		Alias("write", capNetworkWrite, capNetworkRead, capNetworkWrite)

	config := &ResolverConfig{
		Model: model,
		SecretResolver: func(ctx context.Context, s string) (Token, error) {
			if s == "" {
				return &mockToken{false, "", []string{"anonymous"}}, nil
			}
			return repo.FindTokenBySecret(s)
		},
		PolicyResolver: func(ctx context.Context, p string) (Policy, error) {
			return repo.GetPolicyByName(p)
		},
	}

	return config
}

func TestACL(t *testing.T) {

	setupMockRepo()

	t.Run("MissingConfigurations", func(t *testing.T) {
		_, err := NewResolver(&ResolverConfig{})
		if err == nil {
			t.Fatalf("expected error. have %v", err)
		}
	})

	resolver, _ := NewResolver(ACLResolverConfig())

	t.Run("ResolveSecret", func(t *testing.T) {

		secret := "54c06ace-7da6-443b-a5a2-05da5294fbd5"

		acl, err := resolver.ResolveSecret(ctx, secret)
		if err != nil {
			t.Fatal(err)
		}

		if err := acl.CheckAuthorized(ctx, "namespace", "12345", "write-x"); err != nil {
			t.Fatalf("%v", err)
		}

		if err := acl.CheckAuthorized(ctx, "network", "12345", "write"); err == nil {
			t.Fatalf("%v", err)
		}

		if err := acl.CheckAuthorized(ctx, "host", "12345", "write"); err == nil {
			t.Fatalf("%v", err)
		}
	})

	t.Run("InvalidSecret", func(t *testing.T) {

		secret := "foobar"

		_, err := resolver.ResolveSecret(ctx, secret)
		if err == nil {
			t.Fatal(err)
		}
	})

	t.Run("PrivilegedToken", func(t *testing.T) {

		secret := "39076595-19a6-4582-b0d9-bb4a266fd48a"

		acl, err := resolver.ResolveSecret(ctx, secret)
		if err != nil {
			t.Fatal(err)
		}

		if err := acl.CheckAuthorized(ctx, "namespace", "12345", "write-x"); err != nil {
			t.Fatalf("%v", err)
		}

		if err := acl.CheckAuthorized(ctx, "network", "12345", "write"); err != nil {
			t.Fatalf("%v", err)
		}

		if err := acl.CheckAuthorized(ctx, "host", "12345", "write"); err != nil {
			t.Fatalf("%v", err)
		}
	})

	t.Run("InvalidResource", func(t *testing.T) {
		secret := "54c06ace-7da6-443b-a5a2-05da5294fbd5"
		acl, err := resolver.ResolveSecret(ctx, secret)
		if err != nil {
			t.Fatal(err)
		}
		if err := acl.CheckAuthorized(ctx, "foo", "bar", "write-x"); err == nil {
			t.Fatalf("expected error. have %v", err)
		}
	})

	t.Run("AnonymousToken", func(t *testing.T) {

		secret := ""

		acl, err := resolver.ResolveSecret(ctx, secret)
		if err != nil {
			t.Fatal(err)
		}

		if err := acl.CheckAuthorized(ctx, "namespace", "12345", "write-x"); err == nil {
			t.Fatalf("%v", err)
		}

		if err := acl.CheckAuthorized(ctx, "network", "12345", "write"); err == nil {
			t.Fatalf("%v", err)
		}

		if err := acl.CheckAuthorized(ctx, "host", "12345", "list"); err != nil {
			t.Fatalf("%v", err)
		}
	})

	t.Run("StringDump", func(t *testing.T) {

		secret := "54c06ace-7da6-443b-a5a2-05da5294fbd5"

		acl, err := resolver.ResolveSecret(ctx, secret)
		if err != nil {
			t.Fatal(err)
		}

		str := acl.String()
		if str == "" {
			t.Fatal(err)
		}
	})

}
