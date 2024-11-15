package ports

import "github.com/Taluu/media-go/pkg/domain/media/ports/http"

var (
	NewHttpTagsList    = http.NewHttpListServer
	NewHttpTagCreate   = http.NewTagsCreateServer
	NewHttpMediaSeatch = http.NewMediaSearchHTTPPort
	NewHttpMediaCreate = http.NewMediaCreateHTTPServer
	NewHttpMediaViewer = http.NewMediaViewerHTTPServer
)
