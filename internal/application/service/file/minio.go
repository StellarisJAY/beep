package file

import (
	"beep/internal/config"
	"beep/internal/types/interfaces"
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	client *minio.Client
}

func NewMinio(config *config.Config) (interfaces.FileStore, error) {
	client, err := minio.New(config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Minio.AccessKey, config.Minio.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &Minio{
		client: client,
	}, nil
}

func (s *Minio) Upload(ctx context.Context, bucket string, content io.Reader) (string, error) {
	filename := uuid.New().String()
	info, err := s.client.PutObject(ctx, bucket, filename, content, -1, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return info.Key, nil
}

func (s *Minio) Download(ctx context.Context, bucket, key string) (io.Reader, error) {
	info, err := s.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *Minio) Delete(ctx context.Context, bucket, key string) error {
	return s.client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
}

func (s *Minio) TempDownloadURL(ctx context.Context, bucket, key string) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, bucket, key, time.Hour, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *Minio) CreateBucket(ctx context.Context, name string) error {
	err := s.client.MakeBucket(ctx, name, minio.MakeBucketOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Minio) DeleteBucket(ctx context.Context, name string) error {
	objects := s.client.ListObjects(ctx, name, minio.ListObjectsOptions{
		Recursive: true,
	})
	chanErrors := s.client.RemoveObjects(ctx, name, objects, minio.RemoveObjectsOptions{})
	for err := range chanErrors {
		if err.Err != nil {
			return err.Err
		}
	}
	err := s.client.RemoveBucket(ctx, name)
	if err != nil {
		return err
	}
	return nil
}
