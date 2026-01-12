package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ImageStorage interface {
    Upload(ctx context.Context, key string, r io.Reader, size int64, contentType string) (publicURL string, err error)
    Delete(ctx context.Context, key string) error
}

type MinIOStorage struct {
	client  *minio.Client
	bucket  string
	baseURL string
}

func NewMinIOStorage(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinIOStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	baseURL := "http://" + endpoint + "/" + bucket

	return &MinIOStorage{
		client:  client,
		bucket:  bucket,
		baseURL: baseURL,
	}, nil
}

func (m *MinIOStorage) Upload(ctx context.Context, key string, r io.Reader, size int64, contentType string) (string, error) {
	_, err := m.client.PutObject(ctx, m.bucket, key, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", m.baseURL, key), nil
}

func (m *MinIOStorage) Delete(ctx context.Context, key string) error {
	return m.client.RemoveObject(ctx, m.bucket, key, minio.RemoveObjectOptions{})
}
