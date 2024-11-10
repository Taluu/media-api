package media

import (
	"context"
)

type Media struct {
	ID   string
	Name string
}

type MediaRepository interface {
	GetByIDs(ctx context.Context, mediaIDs ...string) (map[string]Media, error)
	Create(ctx context.Context, name string) (Media, error)
}

type MediaService interface {
	SearchByTag(ctx context.Context, tagName string) ([]Media, map[string][]Tag, error)
	Create(ctx context.Context, name string, tags []string) (Media, []Tag, error)
}
