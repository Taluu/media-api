package file

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Taluu/media-go/pkg/domain/media"
)

func NewUploader(dir string) media.MediaUploader {
	return &fileUploader{dir}
}

type fileUploader struct {
	directory string
}

func (u *fileUploader) GetContent(ctx context.Context, id string) (fileContent []byte, err error) {
	fileContent, err = os.ReadFile(fmt.Sprintf("%s/%s", u.directory, id))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = media.FileNotFound(id)
		}

		err = media.FileError(id, err)
	}

	return
}

func (u *fileUploader) Upload(ctx context.Context, id string, fileContent []byte) error {
	return os.WriteFile(fmt.Sprintf("%s/%s", u.directory, id), fileContent, 0644)
}
