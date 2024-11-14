package media

import (
	"context"
)

type Media struct {
	ID       string
	Name     string
	Mimetype string
}

type MediaRepository interface {
	GetByIDs(ctx context.Context, mediaIDs ...string) (map[string]Media, error)
	Create(ctx context.Context, name string, mimetype string) (Media, error)
}

type MediaService interface {
	SearchByTag(ctx context.Context, tagName string) ([]Media, map[string][]Tag, error)
	Create(ctx context.Context, name string, tags []string, fileContent []byte, mimetype string) (Media, []Tag, error)
}

type MediaUploader interface {
	// Upload uploads the media to the storage
	Upload(ctx context.Context, mediaID string, fileContent []byte) error

	// GetContent gets the content for a media, and returns it.
	// A ErrFile will be returned if something goes wrong.
	GetContent(ctx context.Context, mediaID string) (fileContent []byte, err error)
}
