package media

import (
	"context"
)

type Tag struct {
	Name string
}

type TagRegistry interface {
	GetAll(ctx context.Context) (map[string]Tag, error)
	GetMediaIDsForTag(ctx context.Context, name string) ([]string, error)
	GetTagsForMedias(ctx context.Context, mediasID ...string) (map[string][]Tag, error)
	Create(ctx context.Context, name string) (Tag, error)
	Link(ctx context.Context, tagID, mediaID string) error
}

type TagService interface {
	GetAll(ctx context.Context) ([]Tag, error)
	Create(ctx context.Context, name string) (Tag, error)
}
