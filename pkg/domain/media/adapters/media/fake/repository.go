package fake

import (
	"context"
	"sync"

	//lint:ignore ST1001
	. "github.com/Taluu/media-go/pkg/domain/media"
	"github.com/google/uuid"
)

func NewFake() MediaRepository {
	return &repository{
		medias: make(map[string]Media),
	}
}

type repository struct {
	medias map[string]Media
	mtx    sync.RWMutex
}

func (r *repository) Create(ctx context.Context, name string, mimetype string) (Media, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	id := uuid.NewString()
	media := Media{
		ID:       id,
		Name:     name,
		Mimetype: mimetype,
	}

	r.medias[id] = media

	return media, nil
}

func (r *repository) GetByIDs(ctx context.Context, ids ...string) (map[string]Media, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	result := make(map[string]Media, len(ids))
	for _, id := range ids {
		media, exists := r.medias[id]

		if !exists {
			continue
		}

		result[id] = media
	}

	return result, nil
}
