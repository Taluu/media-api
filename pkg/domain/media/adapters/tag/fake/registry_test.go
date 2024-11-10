package fake

import (
	"context"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	repository := NewFake()
	tag, err := repository.Create(ctx, "foo")

	if err != nil {
		t.Fatalf("error while creating tag object : %e", err)
	}

	if tag.Name != "foo" {
		t.Fatalf("tag is not named as expected %q, had %q", "foo", tag.Name)
	}
}

func TestGetAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repository := NewFake()

	// create 2 tags
	repository.Create(ctx, "foo")
	repository.Create(ctx, "bar")

	tags, err := repository.GetAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error : %e", err)
	}

	if len(tags) != 2 {
		t.Fatalf("not all tags returned")
	}
}

func TestGetMediaIDsForTag(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repository := NewFake()

	repository.Link(ctx, "foo", "media-1")
	repository.Link(ctx, "foo", "media-1") // handing duplicates
	repository.Link(ctx, "foo", "media-2")
	repository.Link(ctx, "bar", "media-1")
	repository.Link(ctx, "bar", "media-3")

	medias, err := repository.GetMediaIDsForTag(ctx, "foo")
	if err != nil {
		t.Fatalf("unexpected errors when getting media ids from a tag : %e", err)
	}

	if len(medias) != 2 {
		t.Fatalf("expected 2 medias to be returned, got %d", len(medias))
	}
}

func TestGetTagsForMedias(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repository := NewFake()

	// the purpose for these fixtures is to have 2 medias sharing a tag,
	// and another media having none in common with the other two.
	repository.Link(ctx, "foo", "media-1")
	repository.Link(ctx, "bar", "media-1")
	repository.Link(ctx, "foo", "media-2")
	repository.Link(ctx, "baz", "media-3")

	tags, err := repository.GetTagsForMedias(ctx, "media-1", "media-2")
	if err != nil {
		t.Fatalf("unexpected errors when getting media ids from a tag : %e", err)
	}

	if len(tags["media-1"]) != 2 {
		t.Fatalf("expected 2 tags for the media-1, got %d", len(tags["media-1"]))
	}

	if len(tags["media-1"]) != 2 {
		t.Fatalf("expected 2 tags for the media-1, got %d", len(tags["media-1"]))
	}
}
