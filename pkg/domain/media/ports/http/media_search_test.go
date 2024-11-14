package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Taluu/media-go/pkg/domain/media/adapters"
	"github.com/Taluu/media-go/pkg/domain/media/services"
)

func TestMediaSearch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	repository := adapters.NewFakeMediaRepository()
	tagRegistry := adapters.NewFakeTagRegistry()
	service := services.NewMediaService(repository, tagRegistry, adapters.NewFakeUploader())
	server := NewMediaSearchHTTPPort(service)

	// fixtures
	service.Create(ctx, "media-1", []string{"tag-1", "tag-2"}, nil, "")
	service.Create(ctx, "media-2", []string{"tag-1", "tag-3"}, nil, "")
	service.Create(ctx, "media-3", []string{"tag-2", "tag-3"}, nil, "")

	t.Run("empty tag", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/medias/", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, r)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.Header.Get("Content-type") != "application/json" {
			t.Fatalf("Expected a application/json content-type, got %s", resp.Header.Get("Content-type"))
		}

		if resp.StatusCode != 400 {
			t.Fatalf("Did not expect HTTP %d (%s)", resp.StatusCode, resp.Status)
		}

		var gotResponse httpError

		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&gotResponse)

		if gotResponse.Code != 400 {
			t.Fatalf("expected an error with a code 400, got %d", gotResponse.Code)
		}

		if gotResponse.Error != "empty tag" {
			t.Fatalf("expected an error with a message %q, got %q", "empty tag", gotResponse.Error)
		}
	})

	t.Run("no media", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/medias/tag-4", nil)
		r.SetPathValue("tag", "tag-4")
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

		var gotResponse mediasSearchHTTP

		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&gotResponse)

		if len(gotResponse.Medias) != 0 {
			t.Fatalf("expected to find no media, found %d", len(gotResponse.Medias))
		}
	})

	t.Run("with medias", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/medias/tag-1", nil)
		r.SetPathValue("tag", "tag-1")
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

		var gotResponse mediasSearchHTTP

		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&gotResponse)

		if len(gotResponse.Medias) != 2 {
			t.Fatalf("expected to find 2 medias, found %d", len(gotResponse.Medias))
		}

		for _, m := range gotResponse.Medias {
			if len(m.Tags) != 2 {
				t.Fatalf("expected a media with 2 tags, got %d", len(m.Tags))
			}

			found := false
			for _, t := range m.Tags {
				if t == "tag-1" {
					found = true
				}
			}

			if found == false {
				t.Fatalf("media %q does not have tag %q", m.ID, "tag-1")
			}

			if m.File != fmt.Sprintf("http://%s/viewer/%s", r.Host, m.ID) {
				t.Fatalf("expected the file property to set the url to view the media, got %q", m.File)
			}
		}
	})
}
