package service

import (
	"time"

	"github.com/gosimple/slug"
	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/repository"
)

// ITagService etiket işlemleri için service interface
type ITagService interface {
	ListTags(offset, limit int) ([]*domain.Tag, int64, error)
	GetTagByID(id uint) (*domain.Tag, error)
	GetTagBySlug(slug string) (*domain.Tag, error)
	CreateTag(name, slug string) (*domain.Tag, error)
	UpdateTag(id uint, name, slug string) (*domain.Tag, error)
	DeleteTag(id uint) error
	GetPopularTags(limit int) ([]*domain.Tag, error)
}

// TagService etiket servisinin implementasyonu
type TagService struct {
	tagRepo repository.ITagRepository
}

// NewTagService yeni bir TagService oluşturur
func NewTagService(tagRepo repository.ITagRepository) ITagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

// ListTags etiketleri listeler
func (s *TagService) ListTags(offset, limit int) ([]*domain.Tag, int64, error) {
	return s.tagRepo.List(offset, limit)
}

// GetTagByID ID'ye göre etiket getirir
func (s *TagService) GetTagByID(id uint) (*domain.Tag, error) {
	return s.tagRepo.GetByID(id)
}

// GetTagBySlug slug'a göre etiket getirir
func (s *TagService) GetTagBySlug(slug string) (*domain.Tag, error) {
	return s.tagRepo.GetBySlug(slug)
}

// CreateTag yeni bir etiket oluşturur
func (s *TagService) CreateTag(name, tagSlug string) (*domain.Tag, error) {
	// Slug belirtilmemişse, isimden oluştur
	if tagSlug == "" {
		tagSlug = slug.Make(name)
	}

	now := time.Now()
	tag := &domain.Tag{
		Name:      name,
		Slug:      tagSlug,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// UpdateTag etiketi günceller
func (s *TagService) UpdateTag(id uint, name, tagSlug string) (*domain.Tag, error) {
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		tag.Name = name
	}

	if tagSlug != "" {
		tag.Slug = tagSlug
	} else if name != "" {
		// İsim değiştiyse ve slug belirtilmediyse, isimden yeni slug oluştur
		tag.Slug = slug.Make(name)
	}

	tag.UpdatedAt = time.Now()

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// DeleteTag etiketi siler
func (s *TagService) DeleteTag(id uint) error {
	return s.tagRepo.Delete(id)
}

// GetPopularTags popüler etiketleri getirir
func (s *TagService) GetPopularTags(limit int) ([]*domain.Tag, error) {
	return s.tagRepo.GetPopular(limit)
}
