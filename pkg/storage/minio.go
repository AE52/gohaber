package storage

import (
	"context"
	"io"
	"log"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioConfig MinIO yapılandırma bilgilerini içerir
type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	Location        string
}

// MinioService MinIO işlemleri için servis
type MinioService struct {
	Client     *minio.Client
	BucketName string
	Location   string
}

// NewMinioService yeni bir MinIO servisi oluşturur
func NewMinioService(config *MinioConfig) (*MinioService, error) {
	// MinIO client oluştur
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Bucket kontrolü yap, yoksa oluştur
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.BucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{
			Region: config.Location,
		})
		if err != nil {
			return nil, err
		}
		log.Printf("Bucket oluşturuldu: %s", config.BucketName)
	}

	return &MinioService{
		Client:     client,
		BucketName: config.BucketName,
		Location:   config.Location,
	}, nil
}

// UploadFile dosyayı MinIO'ya yükler
func (s *MinioService) UploadFile(fileReader io.Reader, fileName string, fileSize int64, contentType string) (string, error) {
	ctx := context.Background()

	// Dosya uzantısını kontrol et
	ext := filepath.Ext(fileName)
	if ext == "" {
		// Varsayılan olarak .bin uzantısı ekle
		fileName = fileName + ".bin"
	}

	// Dosya adı prefix'i oluştur (YYYY-MM-DD/)
	prefix := time.Now().Format("2006-01-02") + "/"
	objectName := prefix + fileName

	// Dosyayı yükle
	_, err := s.Client.PutObject(ctx, s.BucketName, objectName, fileReader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

// GetFileURL MinIO'da depolanan dosyanın URL'sini döner
func (s *MinioService) GetFileURL(objectName string, expires time.Duration) (string, error) {
	ctx := context.Background()

	// Geçici URL oluştur
	presignedURL, err := s.Client.PresignedGetObject(ctx, s.BucketName, objectName, expires, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// DeleteFile MinIO'dan dosyayı siler
func (s *MinioService) DeleteFile(objectName string) error {
	ctx := context.Background()

	err := s.Client.RemoveObject(ctx, s.BucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

// ListFiles MinIO'daki dosyaları listeler
func (s *MinioService) ListFiles(prefix string) ([]minio.ObjectInfo, error) {
	ctx := context.Background()

	objects := make([]minio.ObjectInfo, 0)
	objectCh := s.Client.ListObjects(ctx, s.BucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, object)
	}

	return objects, nil
}
