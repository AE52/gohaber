package service

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/repository"
	"github.com/username/haber/pkg/storage"
)

// IMediaService medya işlemleri için servis arayüzü
type IMediaService interface {
	UploadFile(file *multipart.FileHeader, folder string, userID uint) (*domain.Media, error)
	GetFileURL(media *domain.Media, expires time.Duration) (string, error)
	DeleteFile(id uint) error
	GetMediaByID(id uint) (*domain.Media, error)
	GetMediaByUser(userID uint) ([]*domain.Media, error)
	ListMedia(offset, limit int) ([]*domain.Media, int64, error)
	GetMediaContent(objectName string) ([]byte, error)
	DeleteMedia(id uint, userID uint) error
}

// MediaService MediaService'in implementasyonu
type MediaService struct {
	mediaRepo   repository.IMediaRepository
	minioClient *storage.MinioService
}

// NewMediaService yeni bir MediaService oluşturur
func NewMediaService(mediaRepo repository.IMediaRepository, minioClient *storage.MinioService) IMediaService {
	return &MediaService{
		mediaRepo:   mediaRepo,
		minioClient: minioClient,
	}
}

// ErrMinioServiceUnavailable MinIO servisi kullanılamadığında döndürülen hata
var ErrMinioServiceUnavailable = errors.New("MinIO servisi kullanılamıyor, medya işlemleri devre dışı")

// UploadFile dosyayı MinIO'ya yükler ve medya kaydı oluşturur
func (s *MediaService) UploadFile(file *multipart.FileHeader, folder string, userID uint) (*domain.Media, error) {
	// MinIO servisi bağlantısını kontrol et
	if s.minioClient == nil {
		return nil, ErrMinioServiceUnavailable
	}

	// Dosyayı aç
	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	// Dosyayı MinIO'ya yükle
	objectName, err := s.minioClient.UploadFile(
		fileReader,
		file.Filename,
		int64(file.Size),
		file.Header.Get("Content-Type"),
	)
	if err != nil {
		return nil, err
	}

	// Medya bilgilerini veritabanına kaydet
	media := &domain.Media{
		Filename:    filepath.Base(file.Filename),
		ObjectName:  objectName,
		ContentType: file.Header.Get("Content-Type"),
		Filesize:    uint(file.Size),
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Veritabanına kaydet
	if err := s.mediaRepo.Create(media); err != nil {
		// Veritabanına kayıt başarısız olursa, MinIO'dan da sil
		_ = s.minioClient.DeleteFile(objectName) // Hata görmezden geliniyor
		return nil, err
	}

	return media, nil
}

// GetFileURL dosyanın geçici URL'sini oluşturur
func (s *MediaService) GetFileURL(media *domain.Media, expires time.Duration) (string, error) {
	if s.minioClient == nil {
		return "", ErrMinioServiceUnavailable
	}
	return s.minioClient.GetFileURL(media.ObjectName, expires)
}

// DeleteFile dosya ve medya kaydını siler
func (s *MediaService) DeleteFile(id uint) error {
	// Medya kaydını bul
	media, err := s.mediaRepo.Get(id)
	if err != nil {
		return err
	}

	// MinIO servisi bağlantısını kontrol et
	if s.minioClient == nil {
		return ErrMinioServiceUnavailable
	}

	// MinIO'dan dosyayı sil
	err = s.minioClient.DeleteFile(media.ObjectName)
	if err != nil {
		return err
	}

	// DB'den medya kaydını sil
	return s.mediaRepo.Delete(id)
}

// GetMediaByID ID'ye göre medya kaydını bulur
func (s *MediaService) GetMediaByID(id uint) (*domain.Media, error) {
	return s.mediaRepo.Get(id)
}

// GetMediaByUser kullanıcıya ait medya kayıtlarını listeler
func (s *MediaService) GetMediaByUser(userID uint) ([]*domain.Media, error) {
	result, _, err := s.mediaRepo.GetByUser(userID, 0, 1000) // Şimdilik tüm kayıtları getir
	return result, err
}

// ListMedia medya kayıtlarını listeler
func (s *MediaService) ListMedia(offset, limit int) ([]*domain.Media, int64, error) {
	return s.mediaRepo.List(offset, limit)
}

// GetMediaContent medya içeriğini getirir
func (s *MediaService) GetMediaContent(objectName string) ([]byte, error) {
	// Bu fonksiyon henüz implement edilmedi
	return nil, errors.New("not implemented")
}

// DeleteMedia medyayı siler
func (s *MediaService) DeleteMedia(id uint, userID uint) error {
	// Medyayı getir
	media, err := s.mediaRepo.Get(id)
	if err != nil {
		return err
	}

	// Sadece medyayı yükleyen kullanıcı veya admin silebilir
	// Admin kontrolü burada yapılmıyor, üst katmanda yapılmalı
	if media.UserID != userID {
		return domain.ErrForbidden
	}

	// MinIO servisi bağlantısını kontrol et
	if s.minioClient == nil {
		return ErrMinioServiceUnavailable
	}

	// Önce MinIO'dan sil
	if err := s.minioClient.DeleteFile(media.ObjectName); err != nil {
		return err
	}

	// Veritabanından sil
	return s.mediaRepo.Delete(id)
}
