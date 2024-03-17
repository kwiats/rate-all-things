package media

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LocalBlobStorage struct {
	Path string
}

func NewBlobStorage(path string) *LocalBlobStorage {
	return &LocalBlobStorage{
		Path: path,
	}
}

func (storage *LocalBlobStorage) SaveFile(file []byte) (string, error) {
	var fileExt string
	switch {
	case len(file) > 3 && file[0] == 0xFF && file[1] == 0xD8 && file[2] == 0xFF:
		fileExt = ".jpg"
	case len(file) > 3 && file[0] == 0x89 && file[1] == 'P' && file[2] == 'N' && file[3] == 'G':
		fileExt = ".png"
	default:
		return "", errors.New("unsupported image type")
	}

	fileName := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), "image", fileExt)

	if err := os.MkdirAll(storage.Path, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating directory: %w", err)
	}

	fullPath := filepath.Join(storage.Path, fileName)

	if err := os.WriteFile(fullPath, file, os.ModePerm); err != nil {
		return "", fmt.Errorf("error writing file: %w", err)
	}

	return fullPath, nil
}
