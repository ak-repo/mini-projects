package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client *minio.Client
	Bucket string
	UseSSL bool
}

func NewMinioClient() *MinioClient {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")
	useSsl := os.Getenv("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSsl,
	})
	if err != nil {
		log.Fatalf("failed to create minio client: %v", err)
	}

	// ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatalf("failed to check bucket: %v", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("failed to create bucket: %v", err)
		}
		fmt.Printf("created bucket %s\n", bucket)
	}

	return &MinioClient{
		Client: client,
		Bucket: bucket,
		UseSSL: useSsl,
	}
}

// Upload a file from disk (helper)
func (m *MinioClient) FPutObject(ctx context.Context, objectName, filePath, contentType string) (minio.UploadInfo, error) {
	opts := minio.PutObjectOptions{ContentType: contentType}
	return m.Client.FPutObject(ctx, m.Bucket, objectName, filePath, opts)
}

// Upload from Reader (e.g. memory / stream)
func (m *MinioClient) PutObject(ctx context.Context, objectName string, reader interface{}, size int64, contentType string) (minio.UploadInfo, error) {
	opts := minio.PutObjectOptions{ContentType: contentType}
	return m.Client.PutObject(ctx, m.Bucket, objectName, reader.(io.Reader), size, opts)
}

// Generate presigned GET url
func (m *MinioClient) PresignedGetURL(ctx context.Context, objectName string, expirySeconds int64) (string, error) {
	presignedURL, err := m.Client.PresignedGetObject(ctx, m.Bucket, objectName, time.Second*time.Duration(expirySeconds), url.Values{})
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// Convenience: get expiry from env
func SignedURLExpiry() int64 {
	s := os.Getenv("SIGNED_URL_EXPIRY")
	if s == "" {
		return 3600
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 3600
	}
	return v
}
