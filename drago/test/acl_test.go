package test

import (
	"context"
	"fmt"
	"testing"

	drago "github.com/seashell/drago/drago"
	inmem "github.com/seashell/drago/drago/state/inmem"
	structs "github.com/seashell/drago/drago/structs"
	"github.com/seashell/drago/pkg/uuid"
	"github.com/shurcooL/go-goon"
)

type MockAuthHandler struct{}

func (h *MockAuthHandler) Authorize() error {
	return nil
}

var mockConfig = &drago.Config{}

func TestACLBootstrap(t *testing.T) {

	ctx := context.TODO()

	authHandler := &MockAuthHandler{}
	state := inmem.NewStateRepository(nil)
	service := drago.NewACLService(mockConfig, state, authHandler)

	t.Run("Once", func(t *testing.T) {
		var out structs.ACLTokenUpsertResponse
		err := service.BootstrapACL(&structs.ACLBootstrapRequest{}, &out)
		if err != nil {
			t.Fatal(err)
		}
	})

	b.Clear()

	t.Run("Twice", func(t *testing.T) {
		var out structs.ACLTokenUpsertResponse
		service.BootstrapACL(&structs.ACLBootstrapRequest{}, &out)
		err := service.BootstrapACL(&structs.ACLBootstrapRequest{}, &out)
		if err == nil {
			t.Fatal(err)
		}
	})
}

func TestACLTokens(t *testing.T) {

	ctx := context.TODO()

	authHandler := &MockAuthHandler{}
	state := inmem.NewStateRepository(nil)
	service := drago.NewACLService(mockConfig, state, authHandler)

	var out structs.ACLTokenUpsertResponse

	t.Run("Create", func(t *testing.T) {
		_, err := service.UpsertToken(&structs.ACLTokenUpsertRequest{
			Name:     "my-token-1",
			Type:     "management",
			Policies: nil,
		}, &out)
		if err != nil {
			t.Fatal(err)
		}

		_, err = service.UpsertToken(&structs.ACLTokenUpsertRequest{
			Name:     "my-token-2",
			Type:     "foo",
			Policies: nil,
		}, &out)
		if err == nil {
			t.Fatal(err)
		}

		tokens, _ := state.ACLTokens(ctx)
		if len(tokens) != 1 {
			t.Fatal("failed to create token")
		}
	})

	state.Clear()

	t.Run("GetByID", func(t *testing.T) {

		newToken := &structs.ACLToken{
			ID:     uuid.Generate(),
			Secret: uuid.Generate(),
			Name:   "some-token",
			Type:   "management",
		}

		state.UpsertACLToken(ctx, newToken)

		var out structs.SingleACLTokenResponse
		err := service.GetToken(&structs.ACLTokenSpecificRequest{ID: newToken.ID}, &out)
		if err != nil {
			t.Fatalf("failed to get token by id: %v", err)
		}
		if out.ID != *id {
			t.Fatalf("failed to get token by id: %v", err)
		}
	})

	state.Clear()

	t.Run("List", func(t *testing.T) {

		state.UpsertACLToken(ctx, &structs.ACLToken{ID: uuid.Generate(), Name: "foo", Type: "management"})
		state.UpsertACLToken(ctx, &structs.ACLToken{ID: uuid.Generate(), Name: "bar", Type: "client"})

		var out structs.ACLTokenListResponse
		err := service.ListTokens(&structs.ACLTokenListRequest{}, &out)
		if err != nil {
			t.Fatalf("failed to list tokens: %v", err)
		}
		if len(out.Items) != 2 {
			t.Fatalf("failed to create tokens: %v", err)
		}
	})
}

func TestACLAuthorization(t *testing.T) {

	ctx := context.TODO()

	authHandler := &MockAuthHandler{}
	state := inmem.NewStateRepository(nil)
	service := drago.NewACLService(mockConfig, state, authHandler)

	// Create policies
	state.UpsertACLPolicy(ctx, anonymousPolicy())
	state.UpsertACLPolicy(ctx, servicePolicy())
	state.UpsertACLPolicy(ctx, devicePolicy())

	// Create new management token

	t1 := &structs.ACLToken{ID: uuid.Generate(), Name: "admin-1", Type: "management", Secret: uuid.Generate()}
	t2 := &structs.ACLToken{ID: uuid.Generate(), Name: "device-1", Type: "client", Secret: uuid.Generate(), Policies: []string{"device"}}
	t3 := &structs.ACLToken{ID: uuid.Generate(), Name: "service-1", Type: "client", Secret: uuid.Generate(), Policies: []string{"service"}}

	state.UpsertACLToken(ctx, t1)
	state.UpsertACLToken(ctx, t2)
	state.UpsertACLToken(ctx, t3)

	// Get token secret
	token, _ := state.ACLTokenByID(ctx, t2.ID)

	t.Run("Resolve", func(t *testing.T) {
		var out structs.ResolveACLTokenResponse
		err := service.ResolveToken(&structs.ResolveACLTokenRequest{Secret: token.Secret}, &out)
		if err != nil {
			t.Fatalf("failed to resolve token: %v", err)
		}
		if out.ID != t2.ID {
			t.Fatalf("failed to resolve token: retrieved wrong token")
		}
	})

	t.Run("Evaluate Policies", func(t *testing.T) {
		var out structs.ResolveACLTokenResponse
		out, _ := service.ResolveToken(&structs.ACLResolveTokenRequest{Secret: token.Secret}, &out)
		for _, name := range out.Policies {
			p, err := state.ACLPolicyByName(ctx, name)
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

func anonymousPolicy() *structs.ACLPolicy {
	return &structs.ACLPolicy{
		Name: "anonymous",
		Rules: []*structs.ACLPolicyRule{
			{"network", "*", []string{"write"}},
			{"host", "*", []string{"write"}},
		},
	}
}

func servicePolicy() *structs.ACLPolicy {
	return &structs.ACLPolicy{
		Name: "service",
		Rules: []*structs.ACLPolicyRule{
			{"network", "*", []string{"write"}},
			{"host", "*", []string{"write"}},
		},
	}
}

func devicePolicy() *structs.ACLPolicy {
	return &structs.ACLPolicy{
		Name: "device",
		Rules: []*structs.ACLPolicyRule{
			{"network", "*", []string{"write"}},
			{"host", "*", []string{"write"}},
		},
	}
}
