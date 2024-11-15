package http

import (
	"log"
	"net/http"

	//lint:ignore ST1001
	. "github.com/Taluu/media-go/pkg/domain/media"
)

func NewHttpListServer(service TagService) http.Handler {
	return &tagListServer{service}
}

type tagListServer struct {
	TagService
}

func (s *tagListServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tags, err := s.GetAll(ctx)
	if err != nil {
		log.Println("error while getting the tags : ", err)
		jsonError(w, "internal errror", toHttpCode(err))
		return
	}

	tagsHttp := make([]tagListHttp, len(tags))
	for k, tag := range tags {
		tagsHttp[k] = tagListHttp{Name: tag.Name}
	}

	list := tagsListHttp{Tags: tagsHttp}
	jsonResponse(w, list, http.StatusOK)
}

type tagListHttp struct {
	Name string `json:"name"`
}

type tagsListHttp struct {
	Tags []tagListHttp `json:"tags"`
}
