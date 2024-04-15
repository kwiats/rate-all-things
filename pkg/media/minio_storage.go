package media

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioBlobStorage struct {
	Client *minio.Client
	Bucket string
}

func NewMinioBlobStorage(endpoint, accessKeyID, secretAccessKey string, isSecure bool) (*MinioBlobStorage, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: isSecure,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &MinioBlobStorage{
		Client: minioClient,
	}, nil
}

func (bs *MinioBlobStorage) CreateBucket(bucketName string) error {
	ctx := context.Background()
	err := bs.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}
	return nil
}

func (bs *MinioBlobStorage) BucketExists(bucketName string) (bool, error) {
	ctx := context.Background()
	exists, err := bs.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return false, fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	return exists, nil
}

func (bs *MinioBlobStorage) ListBuckets() ([]minio.BucketInfo, error) {
	ctx := context.Background()
	buckets, err := bs.Client.ListBuckets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}
	return buckets, nil
}

func (bs *MinioBlobStorage) ListFiles(bucketName, prefix string) ([]minio.ObjectInfo, error) {
	ctx := context.Background()
	objectCh := bs.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Prefix: prefix, Recursive: true})
	var objects []minio.ObjectInfo
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error on listing objects: %w", object.Err)
		}
		objects = append(objects, object)
	}
	return objects, nil
}

func (bs *MinioBlobStorage) GetFile(bucketName, objectName string) (*url.URL, error) {
	ctx := context.Background()
	reqParams := make(url.Values)
	presignedURL, err := bs.Client.PresignedGetObject(ctx, bucketName, objectName, time.Hour, reqParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get presigned URL for object: %w", err)
	}
	return presignedURL, nil
}

func (bs *MinioBlobStorage) DeleteBucket(bucketName string) error {
	ctx := context.Background()
	err := bs.Client.RemoveBucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}
	return nil
}

func (bs *MinioBlobStorage) DeleteFile(bucketName, objectName string) error {
	ctx := context.Background()
	err := bs.Client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

func (bs *MinioBlobStorage) SaveFile(fileHeader *multipart.FileHeader, direction string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	filePath := fmt.Sprintf("%s/%s", direction, fileHeader.Filename)
	info, err := bs.Client.PutObject(context.Background(), bs.Bucket, filePath, file, fileHeader.Size, minio.PutObjectOptions{ContentType: fileHeader.Header.Get("Content-Type")})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to minio: %w", err)
	}

	fmt.Printf("Successfully uploaded %s of size %d\n", info.Key, info.Size)
	return filePath, nil
}
