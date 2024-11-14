package media

import (
	"context"
	"testing"
	"time"

	"github.com/Taluu/media-go/pkg/domain/media/adapters"
)

func TestSearchByTag(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	fakeTagRegistry := adapters.NewFakeTagRegistry()
	fakeMediaRepository := adapters.NewFakeMediaRepository()

	// link a few medias
	media1, _ := fakeMediaRepository.Create(ctx, "media-1", "random/mime")
	media2, _ := fakeMediaRepository.Create(ctx, "media-2", "random/mime")
	media3, _ := fakeMediaRepository.Create(ctx, "media-3", "random/mime")

	fakeTagRegistry.Link(ctx, "tag-1", media1.ID)
	fakeTagRegistry.Link(ctx, "tag-2", media1.ID)
	fakeTagRegistry.Link(ctx, "tag-1", media2.ID)
	fakeTagRegistry.Link(ctx, "tag-3", media2.ID)
	fakeTagRegistry.Link(ctx, "tag-4", media3.ID)

	service := NewMediaService(fakeMediaRepository, fakeTagRegistry, adapters.NewFakeUploader())
	medias, tags, err := service.SearchByTag(ctx, "tag-1")
	if err != nil {
		t.Fatalf("an error ocurred while fetching data : %s", err)
	}

	if len(medias) != 2 {
		t.Fatalf("expected 2 medias to be returned, had %d", len(medias))
	}

	if len(tags) != 2 {
		t.Fatalf("expected to get a set of tags for 2 medias, had for %d medias", len(tags))
	}

	if _, exists := tags[media1.ID]; !exists {
		t.Fatalf("expected to have tags associated to the media \"media-1\", got none")
	}

	if len(tags[media1.ID]) != 2 {
		t.Fatalf("expected to get a set of 2 tags for the media \"media-1\", had %d tags", len(tags[media1.ID]))
	}

	if _, exists := tags[media2.ID]; !exists {
		t.Fatalf("expected to have tags associated to the media \"media-2\", got none")
	}

	if len(tags[media2.ID]) != 2 {
		t.Fatalf("expected to get a set of 2 tags for the media \"media-2\", had %d tags", len(tags[media2.ID]))
	}

	if _, exists := tags[media3.ID]; exists {
		t.Fatalf("didn't expect to have media-3 returned in the set of found medias")
	}
}

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	fakeTagRegistry := adapters.NewFakeTagRegistry()
	fakeMediaRepository := adapters.NewFakeMediaRepository()
	fakeUploader := adapters.NewFakeUploader()

	fakeTagRegistry.Create(ctx, "tag-1")

	service := NewMediaService(fakeMediaRepository, fakeTagRegistry, fakeUploader)
	media, tags, err := service.Create(ctx, "media-1", []string{"tag-1", "tag-2"}, []byte("content"), "random/mime")

	if err != nil {
		t.Fatalf("an error ocurred while fetching data : %s", err)
	}

	if len(tags) != 2 {
		t.Fatalf("expected 2 tags to be linked, had %d", len(tags))
	}

	if media.Name != "media-1" {
		t.Fatalf("expected the created media to have the name %q, got %q", "media-1", media.Name)
	}

	_, err = fakeUploader.GetContent(ctx, media.ID)
	if err != nil {
		t.Fatalf("an error occurred while uploading file : %s", err)
	}
}
