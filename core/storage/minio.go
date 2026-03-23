package storage

import (
	"fmt"
	"strings"

	"fooder-backend/core/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClient  *minio.Client
	minioEnabled bool
)

func InitMinioClient(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	if !cfg.Minio.Enabled {
		minioClient = nil
		minioEnabled = false
		return nil
	}

	endpoint := strings.TrimSpace(cfg.Minio.Endpoint)
	accessKey := strings.TrimSpace(cfg.Minio.AccessKeyID)
	secretKey := strings.TrimSpace(cfg.Minio.SecretAccessKey)

	if endpoint == "" {
		return fmt.Errorf("missing APP_MINIO_ENDPOINT")
	}
	if accessKey == "" {
		return fmt.Errorf("missing APP_MINIO_ACCESS_KEY_ID")
	}
	if secretKey == "" {
		return fmt.Errorf("missing APP_MINIO_SECRET_ACCESS_KEY")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("init minio client: %w", err)
	}

	minioClient = client
	minioEnabled = true

	return nil
}

func GetMinioClient() *minio.Client {
	return minioClient
}

func IsMinioEnabled() bool {
	return minioEnabled
}
