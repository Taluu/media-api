package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Taluu/media-go/pkg/domain/media"
)

func NewMediaCreateHTTPServer(service media.MediaService) http.Handler {
	return &mediaCreateServer{service}
}

type mediaCreateServer struct {
	service media.MediaService
}

func (m *mediaCreateServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received HTTP %s /medias\n", r.Method)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var request mediaCreateRequest
	data := r.FormValue("data")
	if err := json.Unmarshal([]byte(r.FormValue("data")), &request); data != "" && err != nil && err != io.EOF {
		log.Printf("could not deserialize body into proper json : %s", err)
		jsonError(w, "json error", http.StatusBadRequest)
		return
	}

	fileContent, fileName, mimetype, err := getFile(r)
	if err != nil {
		log.Printf("Problem while fetching file upload : %s", err)
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// use the flename as a name if not provided
	if request.Name == "" {
		request.Name = fileName
	}

	media, tags, err := m.service.Create(ctx, request.Name, request.Tags, fileContent, mimetype)
	if err != nil {
		log.Printf("could not create media : %s", err)
		jsonError(w, "media creation failed", http.StatusInternalServerError)
		return
	}

	tagsListHttp := make([]string, len(tags))
	for k, tag := range tags {
		tagsListHttp[k] = tag.Name
	}

	mediaResponse := mediaCreateResponse{
		ID:   media.ID,
		Name: media.Name,
		File: fmt.Sprintf("http://%s/viewer/%s", r.Host, media.ID),
		Tags: tagsListHttp,
	}
	jsonResponse(w, mediaResponse, http.StatusCreated)
	log.Printf("HTTP POST /medias : 201")
}

func getFile(r *http.Request) (content []byte, filename string, mimetype string, err error) {
	file, header, err := r.FormFile("media")
	if err != nil {
		log.Printf("could not get file : %s", err)
		err = fmt.Errorf("file not found")
		return
	}

	defer file.Close()

	content, err = io.ReadAll(file)
	if err != nil {
		log.Printf("could not read file : %s", err)
		err = fmt.Errorf("file not readable")
	}

	filename = header.Filename
	mimetype = mime.TypeByExtension(filepath.Ext(filename))
	// if for some reason, no mimetype correctly detected, so let's use the generic
	// octet-stream
	if mimetype == "" {
		mimetype = "application/octet-stream"
	}

	return
}

type mediaCreateRequest struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type mediaCreateResponse struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	File string   `json:"file"`
	Tags []string `json:"tags"`
}
