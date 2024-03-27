package media

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type LocalBlobStorage struct {
	Path string
}

func NewBlobStorage(path string) *LocalBlobStorage {
	return &LocalBlobStorage{
		Path: path,
	}
}

func (storage *LocalBlobStorage) SaveFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	filePath := fmt.Sprintf("%s/%s", storage.Path, fileHeader.Filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return filePath, nil
}
