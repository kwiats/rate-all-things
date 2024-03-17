package media

type MediaService struct {
	BlobStorage BlobStorage
}

type BlobStorage interface {
	SaveFile([]byte) (string, error)
}

func NewMediaService(blobStorage BlobStorage) *MediaService {
	return &MediaService{
		BlobStorage: blobStorage,
	}
}

func (storage *MediaService) SaveFile(file []byte) (string, error) {
	return storage.BlobStorage.SaveFile(file)
}
