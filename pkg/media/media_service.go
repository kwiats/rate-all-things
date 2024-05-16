package media

import "mime/multipart"

type MediaService struct {
	BlobStorage IBlobStorage
}

type IBlobStorage interface {
	SaveFile(*multipart.FileHeader, string) (string, error)
}

func NewMediaService(blobStorage IBlobStorage) *MediaService {
	return &MediaService{
		BlobStorage: blobStorage,
	}
}

func (storage *MediaService) SaveFile(file *multipart.FileHeader, direction string) (string, error) {
	return storage.BlobStorage.SaveFile(file, direction)
}
