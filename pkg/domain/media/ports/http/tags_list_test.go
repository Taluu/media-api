package http

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Taluu/media-go/pkg/domain/media/adapters"
	"github.com/Taluu/media-go/pkg/domain/media/services"
)

func TestTagsList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	registry := adapters.NewFakeTagRegistry()
	service := services.NewTagService(registry)
	server := NewHttpListServer(service)

	// create some fixtures
	registry.Create(ctx, "tag-1")
	registry.Create(ctx, "tag-2")

	r := httptest.NewRequest("GET", "/tags", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.Header.Get("Content-type") != "application/json" {
		t.Fatalf("Expected a application/json content-type, got %s", resp.Header.Get("Content-type"))
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Did not expect HTTP %d (%s)", resp.StatusCode, resp.Status)
	}

	var gotResponse tagsListHttp

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&gotResponse)

	if len(gotResponse.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(gotResponse.Tags))
	}
}
