package http

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Taluu/media-go/pkg/domain/media"
)

func NewTagsCreateServer(service media.TagService) http.Handler {
	return &tagCreateServer{service}
}

type tagCreateServer struct {
	service media.TagService
}

func (t *tagCreateServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received HTTP %s /tags\n", r.Method)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var request tagCreateHttp
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&request); err != nil && err != io.EOF {
		log.Printf("could not deserialize body into proper json : %s", err)
		jsonError(w, "json error", 400)
		return
	}

	if request.Name == "" {
		log.Printf("empty tag name")
		jsonError(w, "empty tag name", 400)
		return
	}

	result, err := t.service.Create(ctx, request.Name)
	if err != nil {
		log.Printf("could not create tag : %s", err)
		jsonError(w, "internal error", 500)
		return
	}

	jsonResponse(w, tagCreateHttp{Name: result.Name}, http.StatusCreated)

	log.Println("HTTP POST /tags : 201")
}

type tagCreateHttp struct {
	Name string `json:"name"`
}
