package fake

import (
	"context"
	"testing"
	"time"

	. "github.com/Taluu/media-go/pkg/domain/media"
)

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	repository := NewFake()
	media, err := repository.Create(ctx, "foo", "random/mime")

	if err != nil {
		t.Fatalf("error while creating media object : %e", err)
	}

	if media.Name != "foo" {
		t.Fatalf("media is not named as expected %q, had %q", "foo", media.Name)
	}
}

func TestGetByIDs(t *testing.T) {
	type testCase struct {
		Name   string
		IDs    []string
		Expect func(*testing.T, map[string]Media)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repository := NewFake()

	// create 2 medias with same tag name in common
	media1, _ := repository.Create(ctx, "foo", "random/mime")
	media2, _ := repository.Create(ctx, "bar", "random/mime")

	// create a media with no relation to the other 2
	repository.Create(ctx, "baz", "random/mime")

	cases := []testCase{
		{
			Name: "2 media ids given, 2 returned",
			IDs:  []string{media1.ID, media2.ID},
			Expect: func(t *testing.T, found map[string]Media) {
				if len(found) != 2 {
					t.Fatalf("expected to find 2 medias, found %d", len(found))
				}
			},
		},

		{
			Name: "2 media ids given, 1 found",
			IDs:  []string{media1.ID, "foo"},
			Expect: func(t *testing.T, found map[string]Media) {
				if len(found) != 1 {
					t.Fatalf("expected to find 1 medias, found %d", len(found))
				}
			},
		},

		{
			Name: "none found",
			IDs:  []string{"oops"},
			Expect: func(t *testing.T, found map[string]Media) {
				if len(found) > 0 {
					t.Fatalf("did not expect to find any matching media, found %d", len(found))
				}
			},
		},
	}

	for _, testcase := range cases {
		t.Run(testcase.Name, func(t *testing.T) {
			medias, err := repository.GetByIDs(ctx, testcase.IDs...)
			if err != nil {
				t.Fatalf("unexpected error : %e", err)
			}

			testcase.Expect(t, medias)
		})
	}
}
