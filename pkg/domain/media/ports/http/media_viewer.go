package http

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Taluu/media-go/pkg/domain/media"
)

func NewMediaViewerHTTPServer(service media.MediaService) http.Handler {
	return &mediaViewerServer{service}
}

type mediaViewerServer struct {
	service media.MediaService
}

func (s *mediaViewerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	content, mimetype, err := s.service.View(ctx, r.PathValue("id"))
	if err != nil {
		log.Printf("error while trying to fetch media : %s", err)
		jsonError(w, "media not found", toHttpCode(err))
		return
	}

	sendContent := bytes.NewReader(content)
	sendSize := int64(len(content))

	w.Header().Set("Content-Type", mimetype)
	w.Header().Set("Content-Length", strconv.FormatInt(sendSize, 10))

	w.WriteHeader(http.StatusOK)
	io.CopyN(w, sendContent, sendSize)
}
