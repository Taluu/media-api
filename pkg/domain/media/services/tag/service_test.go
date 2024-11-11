package tag

import (
	"context"
	"testing"
	"time"

	"github.com/Taluu/media-go/pkg/domain/media/adapters"
)

func TestGetAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	registry := adapters.NewFakeTagRegistry()
	service := NewTagService(registry)

	// populate the registry
	registry.Create(ctx, "foo")
	registry.Create(ctx, "bar")

	all, err := service.GetAll(ctx)

	if err != nil {
		t.Fatalf("unexpected error returned by the service : %e", err)
	}

	if len(all) != 2 {
		t.Fatalf("expected to have 2 elements, got %d", len(all))
	}
}
