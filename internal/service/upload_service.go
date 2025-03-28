package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/repository"
)

// IUploadService dosya yükleme için servis arayüzü
type IUploadService interface {
	UploadFile(file *multipart.FileHeader, folder string, userID uint) (*domain.Media, error)
	GetMedia(id uint) (*domain.Media, error)
	DeleteMedia(id uint) error
	ListMedia(offset, limit int) ([]*domain.Media, int64, error)
}

// UploadService dosya yükleme servisi
type UploadService struct {
	mediaRepo repository.IMediaRepository
	uploadDir string
}

// NewUploadService yeni bir UploadService oluşturur
func NewUploadService(mediaRepo repository.IMediaRepository, uploadDir string) IUploadService {
	return &UploadService{
		mediaRepo: mediaRepo,
		uploadDir: uploadDir,
	}
}

// UploadFile dosya yükler
func (s *UploadService) UploadFile(file *multipart.FileHeader, folder string, userID uint) (*domain.Media, error) {
	// Dosya türünü doğrula
	if !s.isAllowedFileType(file.Filename) {
		return nil, &domain.ValidationError{
			Field:   "file",
			Message: "Desteklenmeyen dosya formatı",
		}
	}

	// Dosya boyutunu kontrol et (10MB)
	if file.Size > 10*1024*1024 {
		return nil, &domain.ValidationError{
			Field:   "file",
			Message: "Dosya boyutu en fazla 10MB olabilir",
		}
	}

	// Klasör yolunu oluştur
	uploadPath := filepath.Join(s.uploadDir, folder)
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return nil, err
	}

	// Benzersiz dosya adı oluştur (zaman tabanlı)
	timeNow := time.Now().UnixNano()
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d%s", timeNow, ext)
	objectName := fmt.Sprintf("%s/%s", folder, fileName)
	filePath := filepath.Join(uploadPath, fileName)

	// Dosyayı aç ve içeriği oluştur
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Dosya içeriğini kopyala
	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Medya nesnesini oluştur
	media := &domain.Media{
		Filename:    file.Filename,
		ObjectName:  objectName,
		ContentType: file.Header.Get("Content-Type"),
		Filesize:    uint(file.Size),
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.mediaRepo.Create(media); err != nil {
		// Hata durumunda dosyayı sil
		os.Remove(filePath)
		return nil, err
	}

	return media, nil
}

// GetMedia medya dosyasını getirir
func (s *UploadService) GetMedia(id uint) (*domain.Media, error) {
	return s.mediaRepo.Get(id)
}

// DeleteMedia medya dosyasını siler
func (s *UploadService) DeleteMedia(id uint) error {
	media, err := s.mediaRepo.Get(id)
	if err != nil {
		return err
	}

	// Dosyayı sil
	filePath := filepath.Join(s.uploadDir, media.ObjectName)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Veritabanından sil
	return s.mediaRepo.Delete(id)
}

// ListMedia medya dosyalarını listeler
func (s *UploadService) ListMedia(offset, limit int) ([]*domain.Media, int64, error) {
	return s.mediaRepo.List(offset, limit)
}

// isAllowedFileType dosya türünün desteklenip desteklenmediğini kontrol eder
func (s *UploadService) isAllowedFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	// İzin verilen dosya tipleri
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".zip"}

	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			return true
		}
	}

	return false
}
