package repository

import (
	"errors"

	"github.com/username/haber/internal/domain"
	"gorm.io/gorm"
)

// IMediaRepository medya işlemleri için repository arayüzü
type IMediaRepository interface {
	Create(media *domain.Media) error
	Get(id uint) (*domain.Media, error)
	Delete(id uint) error
	List(offset, limit int) ([]*domain.Media, int64, error)
	GetByUser(userID uint, offset, limit int) ([]*domain.Media, int64, error)
}

// MediaRepository medya repository implementasyonu
type MediaRepository struct {
	db *gorm.DB
}

// NewMediaRepository yeni bir MediaRepository oluşturur
func NewMediaRepository(db *Database) IMediaRepository {
	return &MediaRepository{
		db: db.DB,
	}
}

// Create yeni bir medya kaydı oluşturur
func (r *MediaRepository) Create(media *domain.Media) error {
	return r.db.Create(media).Error
}

// Get ID'ye göre medya getirir
func (r *MediaRepository) Get(id uint) (*domain.Media, error) {
	var media domain.Media
	err := r.db.First(&media, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{
				ResourceType: domain.ResourceMedia,
				ID:           id,
			}
		}
		return nil, err
	}
	return &media, nil
}

// Delete medya kaydını siler
func (r *MediaRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Media{}, id).Error
}

// List medya kayıtlarını listeler
func (r *MediaRepository) List(offset, limit int) ([]*domain.Media, int64, error) {
	var media []*domain.Media
	var count int64

	err := r.db.Model(&domain.Media{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&media).Error
	if err != nil {
		return nil, 0, err
	}

	return media, count, nil
}

// GetByUser kullanıcıya ait medya kayıtlarını getirir
func (r *MediaRepository) GetByUser(userID uint, offset, limit int) ([]*domain.Media, int64, error) {
	var media []*domain.Media
	var count int64

	err := r.db.Model(&domain.Media{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&media).Error
	if err != nil {
		return nil, 0, err
	}

	return media, count, nil
}
