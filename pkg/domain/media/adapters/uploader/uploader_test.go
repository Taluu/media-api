package uploader

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Taluu/media-go/pkg/domain/media"
	"github.com/Taluu/media-go/pkg/domain/media/adapters/uploader/fake"
	"github.com/Taluu/media-go/pkg/domain/media/adapters/uploader/file"
)

func TestFakeUploader(t *testing.T) {
	test(t, fake.NewUploader())
}

func TestFileUploader(t *testing.T) {
	test(t, file.NewUploader("/tmp"))
}

// this test both the upload and the content fetching
func test(t *testing.T, uploader media.MediaUploader) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	t.Run("uploaded file", func(t *testing.T) {
		fileContent := []byte("content")
		err := uploader.Upload(ctx, "media", fileContent)
		if err != nil {
			t.Fatalf("could not upload a file : %s", err)
		}

		content, err := uploader.GetContent(ctx, "media")
		if err != nil {
			t.Fatalf("could not get the file content : %s", err)
		}

		if !bytes.Equal(fileContent, content) {
			t.Fatalf("did not get the right content : %v, expected %v", content, fileContent)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := uploader.GetContent(ctx, "oops")
		if !errors.Is(err, media.ErrFileNotFound) {
			t.Fatalf("expected a FileNotFound error, got %s", err)
		}
	})
}
