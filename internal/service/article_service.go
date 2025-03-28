package service

import (
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/repository"
)

// IArticleService makale işlemleri için service interface
type IArticleService interface {
	CreateArticle(article *domain.CreateArticleRequest, authorID uint) (*domain.Article, error)
	UpdateArticle(id uint, article *domain.UpdateArticleRequest) (*domain.Article, error)
	GetArticleByID(id uint) (*domain.Article, error)
	GetArticleBySlug(slug string) (*domain.Article, error)
	DeleteArticle(id uint) error
	ListArticles(offset, limit int, filters map[string]interface{}) ([]*domain.Article, int64, error)
	GetFeaturedArticles(limit int) ([]*domain.Article, error)
	GetArticlesByCategory(categoryID uint, offset, limit int) ([]*domain.Article, int64, error)
	GetArticlesByTag(tagID uint, offset, limit int) ([]*domain.Article, int64, error)
	GetArticlesByAuthor(authorID uint, offset, limit int) ([]*domain.Article, int64, error)
}

// ArticleService ArticleService'in implementasyonu
type ArticleService struct {
	articleRepo repository.IArticleRepository
	tagRepo     repository.ITagRepository
}

// NewArticleService yeni bir ArticleService oluşturur
func NewArticleService(articleRepo repository.IArticleRepository, tagRepo repository.ITagRepository) IArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
		tagRepo:     tagRepo,
	}
}

// CreateArticle yeni bir makale oluşturur
func (s *ArticleService) CreateArticle(req *domain.CreateArticleRequest, authorID uint) (*domain.Article, error) {
	// Slug oluştur
	slugText := slug.Make(req.Title)

	// Yayınlanma tarihi varsa parse et
	var publishedAt *time.Time
	if req.PublishedAt != "" && req.Status == domain.ArticleStatusPublished {
		t, err := time.Parse(time.RFC3339, req.PublishedAt)
		if err == nil {
			publishedAt = &t
		}
	}

	// Boş değilse özetini al
	summary := req.Summary
	if summary == "" {
		// İçeriğin ilk 150 karakterini al
		summary = strings.TrimSpace(req.Content)
		if len(summary) > 150 {
			summary = summary[:150] + "..."
		}
	}

	// Makale oluştur
	article := &domain.Article{
		Title:         req.Title,
		Slug:          slugText,
		Content:       req.Content,
		Summary:       summary,
		FeaturedImage: req.FeaturedImage,
		AuthorID:      authorID,
		CategoryID:    req.CategoryID,
		Status:        req.Status,
		IsFeatured:    req.IsFeatured,
		PublishedAt:   publishedAt,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Makaleyi kaydet
	err := s.articleRepo.Create(article)
	if err != nil {
		return nil, err
	}

	// TODO: tagRepo tanımlanınca etiketleri ekle

	return article, nil
}

// UpdateArticle makaleyi günceller
func (s *ArticleService) UpdateArticle(id uint, req *domain.UpdateArticleRequest) (*domain.Article, error) {
	// Makaleyi bul
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Alanları güncelle
	if req.Title != "" {
		article.Title = req.Title
		article.Slug = slug.Make(req.Title)
	}

	if req.Content != "" {
		article.Content = req.Content
	}

	if req.Summary != "" {
		article.Summary = req.Summary
	}

	if req.FeaturedImage != "" {
		article.FeaturedImage = req.FeaturedImage
	}

	if req.CategoryID != 0 {
		article.CategoryID = req.CategoryID
	}

	if req.Status != "" {
		article.Status = req.Status
	}

	if req.IsFeatured != nil {
		article.IsFeatured = *req.IsFeatured
	}

	// Yayınlanma tarihini güncelle
	if req.PublishedAt != "" && article.Status == domain.ArticleStatusPublished {
		t, err := time.Parse(time.RFC3339, req.PublishedAt)
		if err == nil {
			article.PublishedAt = &t
		}
	}

	article.UpdatedAt = time.Now()

	// Makaleyi güncelle
	err = s.articleRepo.Update(article)
	if err != nil {
		return nil, err
	}

	// TODO: tagRepo tanımlanınca etiketleri güncelle

	return article, nil
}

// GetArticleByID ID'ye göre makale getirir
func (s *ArticleService) GetArticleByID(id uint) (*domain.Article, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Görüntülenme sayısını artır
	go s.articleRepo.IncrementViewCount(id)

	return article, nil
}

// GetArticleBySlug slug'a göre makale getirir
func (s *ArticleService) GetArticleBySlug(slug string) (*domain.Article, error) {
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Görüntülenme sayısını artır
	go s.articleRepo.IncrementViewCount(article.ID)

	return article, nil
}

// DeleteArticle makaleyi siler
func (s *ArticleService) DeleteArticle(id uint) error {
	return s.articleRepo.Delete(id)
}

// ListArticles makaleleri listeler
func (s *ArticleService) ListArticles(offset, limit int, filters map[string]interface{}) ([]*domain.Article, int64, error) {
	return s.articleRepo.List(offset, limit, filters)
}

// GetFeaturedArticles öne çıkan makaleleri getirir
func (s *ArticleService) GetFeaturedArticles(limit int) ([]*domain.Article, error) {
	return s.articleRepo.GetFeatured(limit)
}

// GetArticlesByCategory kategoriye göre makaleleri getirir
func (s *ArticleService) GetArticlesByCategory(categoryID uint, offset, limit int) ([]*domain.Article, int64, error) {
	return s.articleRepo.GetByCategory(categoryID, offset, limit)
}

// GetArticlesByTag etikete göre makaleleri getirir
func (s *ArticleService) GetArticlesByTag(tagID uint, offset, limit int) ([]*domain.Article, int64, error) {
	return s.articleRepo.GetByTag(tagID, offset, limit)
}

// GetArticlesByAuthor yazara göre makaleleri getirir
func (s *ArticleService) GetArticlesByAuthor(authorID uint, offset, limit int) ([]*domain.Article, int64, error) {
	return s.articleRepo.GetByAuthor(authorID, offset, limit)
}
