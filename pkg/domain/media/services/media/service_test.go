package media

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Taluu/media-go/pkg/domain/media"
	"github.com/Taluu/media-go/pkg/domain/media/adapters"
	"github.com/google/uuid"
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
	service := NewMediaService(fakeMediaRepository, fakeTagRegistry, fakeUploader)

	fakeTagRegistry.Create(ctx, "tag-1")

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

func TestView(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	fakeMediaRepository := adapters.NewFakeMediaRepository()

	service := NewMediaService(
		fakeMediaRepository,
		adapters.NewFakeTagRegistry(),
		adapters.NewFakeUploader(),
	)

	// fixtures
	mediaOK, _, _ := service.Create(ctx, "media-1", nil, []byte("file content"), "random/type")
	mediaNotUploader, _ := fakeMediaRepository.Create(ctx, "media-2", "")

	t.Run("media does not exists", func(t *testing.T) {
		_, _, err := service.View(ctx, uuid.NewString())
		if err == nil {
			t.Error("Expected an error, got non")
		}

		if !errors.Is(err, media.ErrMediaNotFound) {
			t.Errorf("expected a media not found error, got %q", err)
		}
	})

	t.Run("unknown file", func(t *testing.T) {
		_, _, err := service.View(ctx, mediaNotUploader.ID)
		if err == nil {
			t.Error("Expected an error, got non")
			return
		}

		if !errors.Is(err, media.ErrFileNotFound) {
			t.Errorf("expected a not found error, got %s", err)
		}
	})

	t.Run("nominal", func(t *testing.T) {
		content, mimetype, err := service.View(ctx, mediaOK.ID)
		if err != nil {
			t.Errorf("Unexpected error")
		}

		if !bytes.Equal(content, []byte("file content")) {
			t.Errorf("Not the expected content : expected %q, got %q", "file content", string(content))
		}

		if mimetype != "random/type" {
			t.Errorf("Not the expected type : expected %q, got %q", "random/type", mimetype)
		}
	})
}
