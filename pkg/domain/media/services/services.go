package services

import (
	"github.com/Taluu/media-go/pkg/domain/media/services/media"
	"github.com/Taluu/media-go/pkg/domain/media/services/tag"
)

var (
	NewTagService   = tag.NewTagService
	NewMediaService = media.NewMediaService
)
