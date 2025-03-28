package repository

import (
	"errors"

	"github.com/username/haber/internal/domain"
	"gorm.io/gorm"
)

// ITagRepository etiket işlemleri için repository interface
type ITagRepository interface {
	Create(tag *domain.Tag) error
	GetByID(id uint) (*domain.Tag, error)
	GetBySlug(slug string) (*domain.Tag, error)
	Update(tag *domain.Tag) error
	Delete(id uint) error
	List(offset, limit int) ([]*domain.Tag, int64, error)
	GetPopular(limit int) ([]*domain.Tag, error)
}

// TagRepository etiket repository'sinin implementasyonu
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository yeni bir TagRepository oluşturur
func NewTagRepository(db *Database) ITagRepository {
	return &TagRepository{
		db: db.DB,
	}
}

// Create yeni bir etiket oluşturur
func (r *TagRepository) Create(tag *domain.Tag) error {
	return r.db.Create(tag).Error
}

// GetByID ID'ye göre etiket getirir
func (r *TagRepository) GetByID(id uint) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{
				ResourceType: domain.ResourceTag,
				ID:           id,
			}
		}
		return nil, err
	}
	return &tag, nil
}

// GetBySlug slug'a göre etiket getirir
func (r *TagRepository) GetBySlug(slug string) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{
				ResourceType: domain.ResourceTag,
				Slug:         slug,
			}
		}
		return nil, err
	}
	return &tag, nil
}

// Update etiketi günceller
func (r *TagRepository) Update(tag *domain.Tag) error {
	return r.db.Save(tag).Error
}

// Delete etiketi siler
func (r *TagRepository) Delete(id uint) error {
	// Önce bu etiketi kullanan makaleleri kontrol et
	var count int64
	if err := r.db.Table("article_tags").Where("tag_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return &domain.ValidationError{
			Field:   "id",
			Message: "Bu etiket makaleler tarafından kullanılıyor, silinemez",
		}
	}

	return r.db.Delete(&domain.Tag{}, id).Error
}

// List etiketleri listeler
func (r *TagRepository) List(offset, limit int) ([]*domain.Tag, int64, error) {
	var tags []*domain.Tag
	var count int64

	err := r.db.Model(&domain.Tag{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Order("name ASC").Find(&tags).Error
	if err != nil {
		return nil, 0, err
	}

	return tags, count, nil
}

// GetPopular popüler etiketleri getirir (makalelerde en çok kullanılanlar)
func (r *TagRepository) GetPopular(limit int) ([]*domain.Tag, error) {
	var tags []*domain.Tag

	// Bu sorgu, etiket-makale ilişkisine göre etiketleri popülerliklerine göre sıralar
	// GORM kullanarak ilişkiler üzerinden COUNT yapmak karmaşıklaşabilir
	// Bu nedenle raw SQL sorgusu kullanıyoruz
	err := r.db.Raw(`
		SELECT t.* FROM tags t
		LEFT JOIN article_tags at ON t.id = at.tag_id
		GROUP BY t.id
		ORDER BY COUNT(at.article_id) DESC
		LIMIT ?
	`, limit).Scan(&tags).Error

	if err != nil {
		return nil, err
	}

	return tags, nil
}
