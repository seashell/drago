package test

import (
	"context"
	"fmt"
	"testing"

	application "github.com/seashell/drago/drago/application"
	structs "github.com/seashell/drago/drago/application/structs"
	domain "github.com/seashell/drago/drago/domain"
	inmem "github.com/seashell/drago/drago/inmem"
	"github.com/seashell/drago/pkg/uuid"
	"github.com/shurcooL/go-goon"
)

func TestACLBootstrap(t *testing.T) {

	ctx := context.TODO()

	b := inmem.NewBackend()
	sr := inmem.NewACLStateRepositoryAdapter(b)
	pr := inmem.NewACLPolicyRepositoryAdapter(b)
	tr := inmem.NewACLTokenRepositoryAdapter(b)

	service := application.NewACLService(sr, tr, pr)

	t.Run("Once", func(t *testing.T) {
		_, err := service.Bootstrap(ctx, &structs.ACLBootstrapInput{})
		if err != nil {
			t.Fatal(err)
		}
	})

	b.Clear()

	t.Run("Twice", func(t *testing.T) {
		service.Bootstrap(ctx, &structs.ACLBootstrapInput{})
		_, err := service.Bootstrap(ctx, &structs.ACLBootstrapInput{})
		if err == nil {
			t.Fatal(err)
		}
	})
}

func TestACLTokens(t *testing.T) {

	ctx := context.TODO()

	b := inmem.NewBackend()
	tr := inmem.NewACLTokenRepositoryAdapter(b)

	service := application.NewACLTokenService(tr)

	t.Run("Create", func(t *testing.T) {
		_, err := service.Create(context.Background(), &structs.ACLTokenCreateInput{
			Name:     "my-token-1",
			Type:     "management",
			Policies: nil,
		})
		if err != nil {
			t.Fatal(err)
		}

		_, err = service.Create(context.Background(), &structs.ACLTokenCreateInput{
			Name:     "my-token-2",
			Type:     "foo",
			Policies: nil,
		})
		if err == nil {
			t.Fatal(err)
		}

		tokens, _ := tr.FindAll(ctx)
		if len(tokens) != 1 {
			t.Fatal("failed to create token")
		}
	})

	b.Clear()

	t.Run("GetByID", func(t *testing.T) {
		id, _ := tr.Create(ctx, &domain.ACLToken{Name: "some-token", Type: "management"})
		out, err := service.GetByID(ctx, &structs.ACLTokenGetInput{ID: *id})
		if err != nil {
			t.Fatalf("failed to get token by id: %v", err)
		}
		if out.ID != *id {
			t.Fatalf("failed to get token by id: %v", err)
		}
	})

	b.Clear()

	t.Run("List", func(t *testing.T) {
		tr.Create(ctx, &domain.ACLToken{Name: "foo", Type: "management"})
		tr.Create(ctx, &domain.ACLToken{Name: "bar", Type: "client"})
		tokens, err := service.List(ctx, &structs.ACLTokenListInput{})
		if err != nil {
			t.Fatalf("failed to list tokens: %v", err)
		}
		if len(tokens.Items) != 2 {
			t.Fatalf("failed to create tokens: %v", err)
		}
	})
}

func TestACLAuthorization(t *testing.T) {

	ctx := context.TODO()

	b := inmem.NewBackend()
	sr := inmem.NewACLStateRepositoryAdapter(b)
	pr := inmem.NewACLPolicyRepositoryAdapter(b)
	tr := inmem.NewACLTokenRepositoryAdapter(b)

	service := application.NewACLService(sr, tr, pr)

	// Create policies
	pr.Upsert(ctx, anonymousPolicy())
	pr.Upsert(ctx, servicePolicy())
	pr.Upsert(ctx, devicePolicy())

	// Create new management token
	id1, _ := tr.Create(ctx, &domain.ACLToken{Name: "admin-1", Type: "management", Secret: uuid.Generate()})
	id2, _ := tr.Create(ctx, &domain.ACLToken{Name: "device-1", Type: "client", Secret: uuid.Generate(), Policies: []string{"device"}})
	id3, _ := tr.Create(ctx, &domain.ACLToken{Name: "service-1", Type: "client", Secret: uuid.Generate(), Policies: []string{"service"}})

	fmt.Println(id1, id2, id3)

	// Get token secret
	token, _ := tr.GetByID(ctx, *id2)

	dumpBackendState(b)

	t.Run("Resolve", func(t *testing.T) {
		out, err := service.ResolveToken(ctx, &structs.ACLResolveTokenInput{Secret: token.Secret})
		if err != nil {
			t.Fatalf("failed to resolve token: %v", err)
		}
		if out.ID != *id2 {
			t.Fatalf("failed to resolve token: retrieved wrong token")
		}
	})

	t.Run("Evaluate Policies", func(t *testing.T) {
		out, _ := service.ResolveToken(ctx, &structs.ACLResolveTokenInput{Secret: token.Secret})
		// policies := []domain.ACLPolicy{}
		for _, pname := range out.Policies {
			p, err := pr.GetByName(ctx, pname)
			if err != nil {
				t.Fatalf("failed to evaluate policies: %v", err)
			}
			goon.Dump(p)
		}
	})

}

func dumpBackendState(b *inmem.Backend) {
	fmt.Println("\n------------ Backend dump")
	dumpCh := b.Dump()
	for s := range dumpCh {
		fmt.Println(s)
	}
	fmt.Println("------------")
	fmt.Println("")
}

func anonymousPolicy() *domain.ACLPolicy {
	return &domain.ACLPolicy{
		Name: "anonymous",
		NetworkPolicies: []*domain.NetworkPolicy{
			&domain.NetworkPolicy{Target: "*", Capabilities: []string{domain.NetworkCapabilityWrite}},
		},
		HostPolicies: []*domain.HostPolicy{
			&domain.HostPolicy{Target: "*", Capabilities: []string{domain.NetworkCapabilityWrite}},
		},
	}
}

func servicePolicy() *domain.ACLPolicy {
	return &domain.ACLPolicy{
		Name: "service",
		NetworkPolicies: []*domain.NetworkPolicy{
			&domain.NetworkPolicy{Target: "*", Capabilities: []string{domain.NetworkCapabilityWrite}},
		},
		HostPolicies: []*domain.HostPolicy{
			&domain.HostPolicy{Target: "*", Capabilities: []string{domain.NetworkCapabilityWrite}},
		},
	}
}

func devicePolicy() *domain.ACLPolicy {
	return &domain.ACLPolicy{
		Name: "device",
		NetworkPolicies: []*domain.NetworkPolicy{
			&domain.NetworkPolicy{Target: "*", Capabilities: []string{domain.NetworkCapabilityWrite}},
		},
		HostPolicies: []*domain.HostPolicy{
			&domain.HostPolicy{Target: "*", Capabilities: []string{domain.NetworkCapabilityWrite}},
		},
	}
}
