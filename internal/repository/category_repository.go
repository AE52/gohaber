package repository

import (
	"errors"

	"github.com/username/haber/internal/domain"
	"gorm.io/gorm"
)

// ICategoryRepository kategori işlemleri için repository interface
type ICategoryRepository interface {
	Create(category *domain.Category) error
	GetByID(id uint) (*domain.Category, error)
	GetBySlug(slug string) (*domain.Category, error)
	Update(category *domain.Category) error
	Delete(id uint) error
	List(offset, limit int, filters map[string]interface{}) ([]*domain.Category, int64, error)
	GetChildCategories(parentID uint) ([]*domain.Category, error)
}

// CategoryRepository kategori repository'sinin implementasyonu
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository yeni bir CategoryRepository oluşturur
func NewCategoryRepository(db *Database) ICategoryRepository {
	return &CategoryRepository{
		db: db.DB,
	}
}

// Create yeni bir kategori oluşturur
func (r *CategoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

// GetByID ID'ye göre kategori getirir
func (r *CategoryRepository) GetByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Preload("Children").First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{
				ResourceType: "category",
				ID:           id,
			}
		}
		return nil, err
	}
	return &category, nil
}

// GetBySlug slug'a göre kategori getirir
func (r *CategoryRepository) GetBySlug(slug string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Preload("Children").Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{
				ResourceType: "category",
				Slug:         slug,
			}
		}
		return nil, err
	}
	return &category, nil
}

// Update kategoriyi günceller
func (r *CategoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

// Delete kategoriyi siler
func (r *CategoryRepository) Delete(id uint) error {
	// Önce bu kategoriyi kullanan makalelerin kategori bilgisini kontrol et
	var count int64
	if err := r.db.Model(&domain.Article{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return &domain.ValidationError{
			Field:   "id",
			Message: "Bu kategori makaleler tarafından kullanılıyor, silinemez",
		}
	}

	// Önce alt kategorileri kontrol et
	if err := r.db.Model(&domain.Category{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return &domain.ValidationError{
			Field:   "id",
			Message: "Bu kategorinin alt kategorileri var, önce onları silmelisiniz",
		}
	}

	return r.db.Delete(&domain.Category{}, id).Error
}

// List kategorileri listeler
func (r *CategoryRepository) List(offset, limit int, filters map[string]interface{}) ([]*domain.Category, int64, error) {
	var categories []*domain.Category
	var count int64

	query := r.db.Model(&domain.Category{})

	// Filtreleri uygula
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Children").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&categories).Error

	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

// GetChildCategories belirli bir üst kategoriye ait alt kategorileri getirir
func (r *CategoryRepository) GetChildCategories(parentID uint) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.Where("parent_id = ?", parentID).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
