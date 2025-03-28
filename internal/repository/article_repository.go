package repository

import (
	"errors"

	"github.com/username/haber/internal/domain"
	"gorm.io/gorm"
)

// IArticleRepository makale işlemleri için repository interface
type IArticleRepository interface {
	Create(article *domain.Article) error
	GetByID(id uint) (*domain.Article, error)
	GetBySlug(slug string) (*domain.Article, error)
	Update(article *domain.Article) error
	Delete(id uint) error
	List(offset, limit int, filters map[string]interface{}) ([]*domain.Article, int64, error)
	IncrementViewCount(id uint) error
	GetFeatured(limit int) ([]*domain.Article, error)
	GetByCategory(categoryID uint, offset, limit int) ([]*domain.Article, int64, error)
	GetByTag(tagID uint, offset, limit int) ([]*domain.Article, int64, error)
	GetByAuthor(authorID uint, offset, limit int) ([]*domain.Article, int64, error)
}

// ArticleRepository ArticleRepository'nin GORM implementasyonu
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository yeni bir ArticleRepository instance'ı oluşturur
func NewArticleRepository(db *Database) IArticleRepository {
	return &ArticleRepository{db: db.DB}
}

// Create yeni bir makale oluşturur
func (r *ArticleRepository) Create(article *domain.Article) error {
	return r.db.Create(article).Error
}

// GetByID ID'ye göre makale getirir
func (r *ArticleRepository) GetByID(id uint) (*domain.Article, error) {
	var article domain.Article
	err := r.db.Preload("Author").Preload("Category").Preload("Tags").First(&article, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{
				ResourceType: domain.ResourceArticle,
				ID:           id,
			}
		}
		return nil, err
	}
	return &article, nil
}

// GetBySlug slug'a göre makale getirir
func (r *ArticleRepository) GetBySlug(slug string) (*domain.Article, error) {
	var article domain.Article
	err := r.db.Preload("Author").Preload("Category").Preload("Tags").Where("slug = ?", slug).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{
				ResourceType: domain.ResourceArticle,
				ID:           slug,
			}
		}
		return nil, err
	}
	return &article, nil
}

// Update makaleyi günceller
func (r *ArticleRepository) Update(article *domain.Article) error {
	return r.db.Save(article).Error
}

// Delete makaleyi siler (soft delete)
func (r *ArticleRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Article{}, id).Error
}

// List makaleleri listeler
func (r *ArticleRepository) List(offset, limit int, filters map[string]interface{}) ([]*domain.Article, int64, error) {
	var articles []*domain.Article
	var count int64
	query := r.db

	// Filtreleri uygula
	if filters != nil {
		for key, value := range filters {
			query = query.Where(key, value)
		}
	}

	// Toplam sayıyı al
	err := query.Model(&domain.Article{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Makaleleri getir
	err = query.Preload("Author").Preload("Category").
		Offset(offset).Limit(limit).
		Order("published_at DESC").
		Find(&articles).Error

	if err != nil {
		return nil, 0, err
	}

	return articles, count, nil
}

// IncrementViewCount makale görüntülenme sayısını artırır
func (r *ArticleRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&domain.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// GetFeatured öne çıkan makaleleri getirir
func (r *ArticleRepository) GetFeatured(limit int) ([]*domain.Article, error) {
	var articles []*domain.Article
	err := r.db.Preload("Author").Preload("Category").
		Where("is_featured = ?", true).
		Where("status = ?", domain.ArticleStatusPublished).
		Order("published_at DESC").
		Limit(limit).
		Find(&articles).Error

	if err != nil {
		return nil, err
	}

	return articles, nil
}

// GetByCategory kategoriye göre makaleleri getirir
func (r *ArticleRepository) GetByCategory(categoryID uint, offset, limit int) ([]*domain.Article, int64, error) {
	var articles []*domain.Article
	var count int64

	err := r.db.Model(&domain.Article{}).
		Where("category_id = ? AND status = ?", categoryID, domain.ArticleStatusPublished).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Author").Preload("Category").
		Where("category_id = ? AND status = ?", categoryID, domain.ArticleStatusPublished).
		Offset(offset).Limit(limit).
		Order("published_at DESC").
		Find(&articles).Error

	if err != nil {
		return nil, 0, err
	}

	return articles, count, nil
}

// GetByTag etikete göre makaleleri getirir
func (r *ArticleRepository) GetByTag(tagID uint, offset, limit int) ([]*domain.Article, int64, error) {
	var articles []*domain.Article
	var count int64

	// Etikete göre makale sayısını al
	err := r.db.Model(&domain.Article{}).
		Joins("JOIN article_tags ON articles.id = article_tags.article_id").
		Where("article_tags.tag_id = ? AND articles.status = ?", tagID, domain.ArticleStatusPublished).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Etikete göre makaleleri getir
	err = r.db.Preload("Author").Preload("Category").Preload("Tags").
		Joins("JOIN article_tags ON articles.id = article_tags.article_id").
		Where("article_tags.tag_id = ? AND articles.status = ?", tagID, domain.ArticleStatusPublished).
		Offset(offset).Limit(limit).
		Order("articles.published_at DESC").
		Find(&articles).Error

	if err != nil {
		return nil, 0, err
	}

	return articles, count, nil
}

// GetByAuthor yazara göre makaleleri getirir
func (r *ArticleRepository) GetByAuthor(authorID uint, offset, limit int) ([]*domain.Article, int64, error) {
	var articles []*domain.Article
	var count int64

	err := r.db.Model(&domain.Article{}).
		Where("author_id = ? AND status = ?", authorID, domain.ArticleStatusPublished).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Author").Preload("Category").
		Where("author_id = ? AND status = ?", authorID, domain.ArticleStatusPublished).
		Offset(offset).Limit(limit).
		Order("published_at DESC").
		Find(&articles).Error

	if err != nil {
		return nil, 0, err
	}

	return articles, count, nil
}
