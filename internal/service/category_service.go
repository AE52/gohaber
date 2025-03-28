package service

import (
	"time"

	"github.com/gosimple/slug"
	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/repository"
)

// ICategoryService kategori işlemleri için service interface
type ICategoryService interface {
	ListCategories(offset, limit int, filters map[string]interface{}) ([]*domain.Category, int64, error)
	GetCategoryByID(id uint) (*domain.Category, error)
	GetCategoryBySlug(slug string) (*domain.Category, error)
	CreateCategory(req *domain.CreateCategoryRequest) (*domain.Category, error)
	UpdateCategory(id uint, req *domain.UpdateCategoryRequest) (*domain.Category, error)
	DeleteCategory(id uint) error
}

// CategoryService kategori servisinin implementasyonu
type CategoryService struct {
	categoryRepo repository.ICategoryRepository
}

// NewCategoryService yeni bir CategoryService oluşturur
func NewCategoryService(categoryRepo repository.ICategoryRepository) ICategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// ListCategories kategorileri listeler
func (s *CategoryService) ListCategories(offset, limit int, filters map[string]interface{}) ([]*domain.Category, int64, error) {
	return s.categoryRepo.List(offset, limit, filters)
}

// GetCategoryByID ID'ye göre kategori getirir
func (s *CategoryService) GetCategoryByID(id uint) (*domain.Category, error) {
	return s.categoryRepo.GetByID(id)
}

// GetCategoryBySlug slug'a göre kategori getirir
func (s *CategoryService) GetCategoryBySlug(slug string) (*domain.Category, error) {
	return s.categoryRepo.GetBySlug(slug)
}

// CreateCategory yeni bir kategori oluşturur
func (s *CategoryService) CreateCategory(req *domain.CreateCategoryRequest) (*domain.Category, error) {
	// Slug belirtilmemişse, isimden oluştur
	if req.Slug == "" {
		req.Slug = slug.Make(req.Name)
	}

	now := time.Now()
	category := &domain.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory kategoriyi günceller
func (s *CategoryService) UpdateCategory(id uint, req *domain.UpdateCategoryRequest) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}

	if req.Slug != "" {
		category.Slug = req.Slug
	} else if req.Name != "" {
		// İsim değiştiyse ve slug belirtilmediyse, isimden yeni slug oluştur
		category.Slug = slug.Make(req.Name)
	}

	if req.Description != "" {
		category.Description = req.Description
	}

	category.ParentID = req.ParentID
	category.UpdatedAt = time.Now()

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory kategoriyi siler
func (s *CategoryService) DeleteCategory(id uint) error {
	return s.categoryRepo.Delete(id)
}
