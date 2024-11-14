package media

import "fmt"

var (
	ErrFileNotFound = fmt.Errorf("file not found")
	ErrFile         = fmt.Errorf("file error")
)

func FileNotFound(id string) error {
	return fmt.Errorf("%w : media %q", ErrFileNotFound, id)
}

func FileError(id string, err error) error {
	return fmt.Errorf("%w : %w", ErrFile, err)
}
