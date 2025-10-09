package minio

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/musicman-backend/config"
)

func InitMinioClient(minioConfig config.MinioConfig) (*minio.Client, error) {
	client, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:      credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
		Secure:     minioConfig.UseSSL,
		MaxRetries: minioConfig.MaxRetries,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return client, nil
}

type Minio struct {
	client *minio.Client
}

func NewMinio(client *minio.Client) *Minio {
	return &Minio{
		client: client,
	}
}

func (m *Minio) UploadFile(ctx context.Context, bucketName string, objectName string, filePath string) error {
	_, err := m.client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

func (m *Minio) DownloadFile(ctx context.Context, bucketName string, objectName string, filePath string) error {
	err := m.client.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})

	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	return nil
}

func (m *Minio) GetFileURL(ctx context.Context, bucketName string, objectName string) (string, error) {
	// Generate presigned URL valid for 1 hour
	url, err := m.client.PresignedGetObject(ctx, bucketName, objectName, time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

func (m *Minio) DeleteFile(ctx context.Context, bucketName string, objectName string) error {
	err := m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (m *Minio) CreateBucketIfNotExists(ctx context.Context, bucketName string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if exists {
		return nil
	}

	err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	return nil
}
