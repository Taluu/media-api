package fake

import (
	"context"
	"sync"

	//lint:ignore ST1001
	. "github.com/Taluu/media-go/pkg/domain/media"
)

func NewFake() TagRegistry {
	return &repository{
		tags:   make(map[string][]string),
		medias: make(map[string][]string),
	}
}

type repository struct {
	tags   map[string][]string
	medias map[string][]string
	mtx    sync.RWMutex
}

func (r *repository) GetTagsForMedias(ctx context.Context, mediasID ...string) (map[string][]Tag, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	tags := make(map[string][]Tag)

	for _, mediaID := range mediasID {
		tags[mediaID] = make([]Tag, 0, len(r.medias[mediaID]))

		for _, tag := range r.medias[mediaID] {
			tags[mediaID] = append(tags[mediaID], Tag{Name: tag})
		}
	}

	return tags, nil
}

func (r *repository) GetMediaIDsForTag(ctx context.Context, name string) ([]string, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return r.tags[name], nil
}

func (r *repository) GetAll(ctx context.Context) (map[string]Tag, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	result := make(map[string]Tag, len(r.tags))
	for tag := range r.tags {
		result[tag] = Tag{Name: tag}
	}
	return result, nil
}

func (r *repository) Create(ctx context.Context, name string) (Tag, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	tag := Tag{
		Name: name,
	}

	r.tags[name] = make([]string, 0)

	return tag, nil
}

func (r *repository) Link(ctx context.Context, tagID, mediaID string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, exists := r.medias[mediaID]; !exists {
		r.medias[mediaID] = make([]string, 0)
	}

	// ensure uniqness, no need to have the same tag serveral time for a single
	// media
	for _, mediaTag := range r.medias[mediaID] {
		if mediaTag == tagID {
			return nil
		}
	}
	r.medias[mediaID] = append(r.medias[mediaID], tagID)

	if _, exists := r.tags[tagID]; !exists {
		r.tags[tagID] = make([]string, 0)
	}
	r.tags[tagID] = append(r.tags[tagID], mediaID)
	return nil
}
