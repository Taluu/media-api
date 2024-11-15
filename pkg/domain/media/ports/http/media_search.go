package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Taluu/media-go/pkg/domain/media"
)

func NewMediaSearchHTTPPort(service media.MediaService) http.Handler {
	return &mediaSearchServer{service}
}

type mediaSearchServer struct {
	service media.MediaService
}

func (m *mediaSearchServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tag := r.PathValue("tag")
	if tag == "" {
		log.Println("empty tag")
		jsonError(w, "empty tag", http.StatusBadRequest)
		return
	}

	medias, tags, err := m.service.SearchByTag(ctx, tag)
	if err != nil {
		log.Println("error while getting the medias : ", err)
		jsonError(w, "internal errror", toHttpCode(err))
		return
	}

	mediasHttp := make([]mediaSearchHttp, len(medias))
	for k, media := range medias {
		tagsMedia := make([]string, len(tags[media.ID]))
		for kTag, tag := range tags[media.ID] {
			tagsMedia[kTag] = tag.Name
		}

		mediasHttp[k] = mediaSearchHttp{
			ID:   media.ID,
			Name: media.Name,
			Tags: tagsMedia,
			File: fmt.Sprintf("http://%s/viewer/%s", r.Host, media.ID),
		}
	}

	list := mediasSearchHTTP{Medias: mediasHttp}
	jsonResponse(w, list, http.StatusOK)
}

type mediasSearchHTTP struct {
	Medias []mediaSearchHttp `json:"medias"`
}

type mediaSearchHttp struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	File string   `json:"file"`
	Tags []string `json:"tags"`
}
