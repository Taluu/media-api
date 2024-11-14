package fake

import (
	"context"

	//lint:ignore ST1001 it's the domain
	. "github.com/Taluu/media-go/pkg/domain/media"
)

func NewUploader() MediaUploader {
	return &fakeUploader{
		files: make(map[string][]byte),
	}
}

// Uploads a file in memory rather than disk
type fakeUploader struct {
	files map[string][]byte
}

func (u *fakeUploader) GetContent(ctx context.Context, mediaID string) (fileContent []byte, err error) {
	fileContent, exists := u.files[mediaID]
	if !exists {
		err = FileError(mediaID, FileNotFound(mediaID))
		return
	}

	return fileContent, nil
}

func (u *fakeUploader) Upload(ctx context.Context, id string, fileContent []byte) error {
	contentCopy := make([]byte, len(fileContent))
	copy(contentCopy, fileContent)

	u.files[id] = contentCopy
	return nil
}
