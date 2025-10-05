package minio

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	client     *minio.Client
	bucketName string
}

func NewMinio(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*Minio, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	repo := &Minio{
		client:     client,
		bucketName: bucketName,
	}

	// Create bucket if not exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}

		// Set bucket policy if needed
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + bucketName + `/*"]
				}
			]
		}`

		err = client.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			log.Printf("Warning: Failed to set bucket policy: %v", err)
		}
	}

	return repo, nil
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
