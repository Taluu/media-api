package tag

import (
	//lint:ignore ST1001
	"context"
	"maps"
	"slices"

	. "github.com/Taluu/media-go/pkg/domain/media"
)

func NewTagService(registry TagRegistry) TagService {
	return &service{registry}
}

type service struct {
	TagRegistry
}

func (s *service) GetAll(ctx context.Context) ([]Tag, error) {
	tags, err := s.TagRegistry.GetAll(ctx)
	tagsSlice := slices.Collect(maps.Values(tags))
	return tagsSlice, err
}
