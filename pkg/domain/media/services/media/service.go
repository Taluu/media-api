package media

import (
	"context"
	"maps"
	"slices"

	. "github.com/Taluu/media-go/pkg/domain/media"
	"golang.org/x/sync/errgroup"
)

func NewMediaService(repository MediaRepository, tagRegistry TagRegistry, uploader MediaUploader) MediaService {
	return &service{repository, tagRegistry, uploader}
}

type service struct {
	MediaRepository
	tags     TagRegistry
	uploader MediaUploader
}

// Create implements media.MediaService.
// Subtle: this method shadows the method (MediaRepository).Create of service.MediaRepository.
func (s *service) Create(ctx context.Context, name string, tags []string, fileContent []byte, mimetype string) (Media, []Tag, error) {
	media, err := s.MediaRepository.Create(ctx, name, mimetype)
	if err != nil {
		return Media{}, nil, err
	}

	tagsSlice := make([]Tag, 0, len(tags))

	for _, tag := range tags {
		// silently ignores if this fails
		// the rationale behind this is "tags are not that important for medias,
		// so it's okay if it doesn't add them" and also "if it doesn't exist,
		// then let's create it"
		if err := s.tags.Link(ctx, tag, media.ID); err == nil {
			tagsSlice = append(tagsSlice, Tag{Name: tag})
		}
	}

	err = s.uploader.Upload(ctx, media.ID, fileContent)

	return media, tagsSlice, err
}

func (s *service) SearchByTag(ctx context.Context, tagName string) ([]Media, map[string][]Tag, error) {
	mediaIds, err := s.tags.GetMediaIDsForTag(ctx, tagName)
	if err != nil {
		return nil, nil, err
	}

	mediasSlice := make([]Media, 0)
	tags := make(map[string][]Tag)

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() (err error) {
		medias, err := s.GetByIDs(ctx, mediaIds...)
		mediasSlice = slices.Collect(maps.Values(medias))
		return
	})

	group.Go(func() (err error) {
		tags, err = s.tags.GetTagsForMedias(ctx, mediaIds...)
		return
	})

	err = group.Wait()

	return mediasSlice, tags, err
}
