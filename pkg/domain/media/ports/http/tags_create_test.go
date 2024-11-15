package http

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Taluu/media-go/pkg/domain/media/adapters"
	"github.com/Taluu/media-go/pkg/domain/media/services"
)

func TestTagCreateServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	registry := adapters.NewFakeTagRegistry()
	service := services.NewTagService(registry)
	server := NewTagsCreateServer(service)

	// create some fixtures
	registry.Create(ctx, "tag-1")

	t.Run("failures", func(t *testing.T) {
		type testCase struct {
			name string
			body json.RawMessage

			expectedCode    int
			expectedMessage string
		}

		testCases := []testCase{
			{
				name:            "invalid json",
				body:            json.RawMessage("not a valid json"),
				expectedCode:    400,
				expectedMessage: "json error",
			},

			{
				name:            "empty tag name",
				body:            json.RawMessage(`{"name": ""}`),
				expectedCode:    400,
				expectedMessage: "empty tag name",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				r := httptest.NewRequest("POST", "/tags", strings.NewReader(string(tc.body))).WithContext(ctx)
				w := httptest.NewRecorder()
				server.ServeHTTP(w, r)

				resp := w.Result()
				defer resp.Body.Close()

				if resp.Header.Get("Content-type") != "application/json" {
					t.Fatalf("Expected a application/json content-type, got %s", resp.Header.Get("Content-type"))
				}

				if resp.StatusCode != tc.expectedCode {
					t.Fatalf("Did not expect HTTP %d (%s)", resp.StatusCode, resp.Status)
				}

				var gotResponse httpError

				decoder := json.NewDecoder(resp.Body)
				decoder.Decode(&gotResponse)

				if gotResponse.Code != tc.expectedCode {
					t.Fatalf("expected %d code, got %d", tc.expectedCode, gotResponse.Code)
				}

				if gotResponse.Error != tc.expectedMessage {
					t.Fatalf("expected message %q, got %q", tc.expectedMessage, gotResponse.Error)
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/tags", strings.NewReader(`{"name": "test"}`)).WithContext(ctx)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, r)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.Header.Get("Content-type") != "application/json" {
			t.Fatalf("Expected a application/json content-type, got %s", resp.Header.Get("Content-type"))
		}

		if resp.StatusCode != 201 {
			t.Fatalf("Did not expect HTTP %d (%s)", resp.StatusCode, resp.Status)
		}

		var gotResponse tagCreateHttp

		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&gotResponse)

		if gotResponse.Name != "test" {
			t.Fatalf("expected %q tag, got %q", "test", gotResponse.Name)
		}
	})
}
