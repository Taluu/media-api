package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Taluu/media-go/pkg/domain/media/adapters"
	"github.com/Taluu/media-go/pkg/domain/media/services"
)

func TestMediaCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mediaRepository := adapters.NewFakeMediaRepository()
	tagRegistry := adapters.NewFakeTagRegistry()
	fakeUploader := adapters.NewFakeUploader()
	service := services.NewMediaService(mediaRepository, tagRegistry, fakeUploader)
	server := NewMediaCreateHTTPServer(service)

	testCases := []struct {
		name     string
		data     string
		ext      string
		withFile bool
		asserter func(req *http.Request, resp *http.Response)
	}{
		{
			name: "invalid json body",
			data: "{not valid}",
			asserter: func(req *http.Request, resp *http.Response) {
				if resp.StatusCode != 400 {
					t.Errorf("expected a 400, got %d", resp.StatusCode)
					return
				}

				var gotResponse httpError
				decoder := json.NewDecoder(resp.Body)
				decoder.Decode(&gotResponse)

				if gotResponse.Error != "json error" {
					t.Errorf("expected a error message %q, got %q", "json error", gotResponse.Error)
					return
				}
			},
		},
		{
			name: "missing file",
			asserter: func(req *http.Request, resp *http.Response) {
				if resp.StatusCode != 400 {
					t.Errorf("expected a 400, got %d", resp.StatusCode)
					return
				}

				var gotResponse httpError
				decoder := json.NewDecoder(resp.Body)
				decoder.Decode(&gotResponse)

				if gotResponse.Error != "file not found" {
					t.Errorf("expected a error message %q, got %q", "file not found", gotResponse.Error)
					return
				}
			},
		},
		{
			name:     "unreadable file",
			withFile: true,
			asserter: func(req *http.Request, resp *http.Response) {
				// TODO : find a way to make an unreadble file
			},
		},
		{
			name:     "media creation error",
			withFile: true,
			asserter: func(req *http.Request, resp *http.Response) {
				// TODO : find a way have the media creation fail
				// e.g "faking" the service so it fails on calls
			},
		},
		{
			name:     "no mimetype found",
			data:     `{"name": "foo"}`,
			withFile: true,
			asserter: func(req *http.Request, resp *http.Response) {
				if resp.StatusCode != 201 {
					t.Errorf("expected a 201, got %d", resp.StatusCode)
					return
				}

				var gotResponse mediaCreateResponse
				decoder := json.NewDecoder(resp.Body)
				decoder.Decode(&gotResponse)

				medias, _ := mediaRepository.GetByIDs(ctx, gotResponse.ID)
				mimetype := medias[gotResponse.ID].Mimetype

				if mimetype != "application/octet-stream" {
					t.Errorf("expected a %q mimetype, got %q", "application/octet-stream", mimetype)
					return
				}

				if len(gotResponse.Tags) != 0 {
					t.Errorf("expected a media without any tags, got %d", len(gotResponse.Tags))
					return
				}

				expectedFile := fmt.Sprintf("http://%s/viewer/%s", req.Host, gotResponse.ID)
				if gotResponse.File != expectedFile {
					t.Errorf("expected the response to expose a file with %q, got %q", expectedFile, gotResponse.File)
					return
				}
			},
		},
		{
			name:     "without a name",
			ext:      ".png",
			withFile: true,
			asserter: func(req *http.Request, resp *http.Response) {
				if resp.StatusCode != 201 {
					t.Errorf("expected a 400, got %d", resp.StatusCode)
				}

				var gotResponse mediaCreateResponse
				decoder := json.NewDecoder(resp.Body)
				decoder.Decode(&gotResponse)

				medias, _ := mediaRepository.GetByIDs(ctx, gotResponse.ID)
				mimetype := medias[gotResponse.ID].Mimetype

				if gotResponse.Name != "fixture.png" {
					t.Errorf("expected to get a media with the filename as its name, got %q instead of %q", gotResponse.Name, "fixture.png")
					return
				}

				if mimetype != "image/png" {
					t.Errorf("expected to get a media with a %q mimetype, got %q", "image/png", mimetype)
					return
				}

				if len(gotResponse.Tags) != 0 {
					t.Errorf("expected a media without any tags, got %d", len(gotResponse.Tags))
					return
				}

				expectedFile := fmt.Sprintf("http://%s/viewer/%s", req.Host, gotResponse.ID)
				if gotResponse.File != expectedFile {
					t.Errorf("expected the response to expose a file with %q, got %q", expectedFile, gotResponse.File)
					return
				}
			},
		},
		{
			name:     "nominal case",
			ext:      ".png",
			data:     `{"name": "a horse with no name", "tags": ["musical reference"]}`,
			withFile: true,
			asserter: func(req *http.Request, resp *http.Response) {
				if resp.StatusCode != 201 {
					t.Errorf("expected a 400, got %d", resp.StatusCode)
					return
				}

				var gotResponse mediaCreateResponse
				decoder := json.NewDecoder(resp.Body)
				decoder.Decode(&gotResponse)

				medias, _ := mediaRepository.GetByIDs(ctx, gotResponse.ID)
				mimetype := medias[gotResponse.ID].Mimetype

				if gotResponse.Name != "a horse with no name" {
					t.Errorf("expected to get a media with the filename as its name, got %q instead of %q", gotResponse.Name, "horse with no name")
					return
				}

				if mimetype != "image/png" {
					t.Errorf("expected to get a media with a %q mimetype, got %q", "image/png", mimetype)
					return
				}

				if len(gotResponse.Tags) != 1 {
					t.Errorf("expected a media with one tag, got %d", len(gotResponse.Tags))
					return
				}

				if gotResponse.Tags[0] != "musical reference" {
					t.Errorf("expected the media to be tagged with %q, got %q", "musical reference", gotResponse.Tags[0])
					return
				}

				expectedFile := fmt.Sprintf("http://%s/viewer/%s", req.Host, gotResponse.ID)
				if gotResponse.File != expectedFile {
					t.Errorf("expected the response to expose a file with %q, got %q", expectedFile, gotResponse.File)
					return
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := prepareRequest(tc.data, tc.ext, tc.withFile)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.Header.Get("content-type") != "application/json" {
				t.Errorf("expected a %q content-type, got %q", "application/json", r.Header.Get("content-type"))
				return
			}

			tc.asserter(r, resp)
		})
	}
}

func prepareRequest(data string, ext string, attachFile bool) *http.Request {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	w.WriteField("data", data)

	if attachFile {
		// add a file attachement, the content-type of the file won't matter as we
		// are just interested in the extension of the file.
		// As it's not within the scope of this test, I choose not to bother with
		// checking the real content-type or limit by size, but we definitely should
		// on a prod environment.
		p, _ := w.CreateFormFile("media", fmt.Sprintf("fixture%s", ext))
		p.Write([]byte("sample fixture test"))
	}

	w.Close()

	request := httptest.NewRequest("POST", "/medias", body)
	request.Header.Add("Content-Type", w.FormDataContentType())

	return request
}
