package adapters

import (
	//lint:ignore ST1001
	. "github.com/Taluu/media-go/pkg/domain/media"
	mediaFake "github.com/Taluu/media-go/pkg/domain/media/adapters/media/fake"
	tagFake "github.com/Taluu/media-go/pkg/domain/media/adapters/tag/fake"
)

func NewFakeMediaRepository() MediaRepository {
	return mediaFake.NewFake()
}

func NewFakeTagRegistry() TagRegistry {
	return tagFake.NewFake()
}
