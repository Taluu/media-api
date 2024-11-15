package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Taluu/media-go/pkg/domain/media/adapters"
	"github.com/Taluu/media-go/pkg/domain/media/services"
	"github.com/google/uuid"
)

func TestMediaViewer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mediaRepository := adapters.NewFakeMediaRepository()
	service := services.NewMediaService(
		mediaRepository,
		adapters.NewFakeTagRegistry(),
		adapters.NewFakeUploader(),
	)

	server := NewMediaViewerHTTPServer(service)

	t.Run("media not found", func(t *testing.T) {
		id := uuid.NewString()
		r := httptest.NewRequest("GET", fmt.Sprintf("/medias/%s", id), nil).WithContext(ctx)
		r.SetPathValue("id", id)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, r)
		resp := w.Result()

		if resp.Header.Get("content-type") != "application/json" {
			t.Errorf("Expected a %q content-type, got %q", "application/json", resp.Header.Get("content-type"))
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected a status not found, got %d", resp.StatusCode)
		}

		var gotResponse httpError
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&gotResponse)

		if gotResponse.Code != http.StatusNotFound {
			t.Errorf("Expected a not found result, got %q", gotResponse.Code)
		}

		if gotResponse.Error != "media not found" {
			t.Errorf("expected an error %q, got %q", "media not found", gotResponse.Error)
		}

	})
	t.Run("file not available", func(t *testing.T) {
		mediaOK, _ := mediaRepository.Create(ctx, "my-media", "text/plain")

		r := httptest.NewRequest("GET", fmt.Sprintf("/medias/%s", mediaOK.ID), nil).WithContext(ctx)
		r.SetPathValue("id", mediaOK.ID)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, r)
		resp := w.Result()

		if resp.Header.Get("content-type") != "application/json" {
			t.Errorf("Expected a %q content-type, got %q", "application/json", resp.Header.Get("content-type"))
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected a status not found, got %d", resp.StatusCode)
		}

		var gotResponse httpError
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&gotResponse)

		if gotResponse.Code != http.StatusNotFound {
			t.Errorf("Expected a not found result, got %q", gotResponse.Code)
		}

		if gotResponse.Error != "media not found" {
			t.Errorf("expected an error %q, got %q", "media not found", gotResponse.Error)
		}
	})
	t.Run("nominal", func(t *testing.T) {
		mediaOK, _, _ := service.Create(ctx, "my-media", nil, []byte("file content"), "text/plain")

		r := httptest.NewRequest("GET", fmt.Sprintf("/medias/%s", mediaOK.ID), nil).WithContext(ctx)
		r.SetPathValue("id", mediaOK.ID)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, r)
		resp := w.Result()

		if resp.Header.Get("content-type") != "text/plain" {
			t.Errorf("Expected a %q content-type, got %q", "text/plain", resp.Header.Get("content-type"))
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected a status ok, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		if !bytes.Equal(body, []byte("file content")) {
			t.Errorf("expected %q as file content, got %q", "file content", string(body))
		}
	})
}
