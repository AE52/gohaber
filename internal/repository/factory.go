package repository

import (
	"sync"
)

// RepositoryFactory tüm repository'leri yönetir
type RepositoryFactory struct {
	db           *Database
	userRepo     IUserRepository
	articleRepo  IArticleRepository
	categoryRepo ICategoryRepository
	tagRepo      ITagRepository
	mediaRepo    IMediaRepository
	mu           sync.RWMutex
}

// NewRepositoryFactory yeni bir factory oluşturur
func NewRepositoryFactory(db *Database) *RepositoryFactory {
	return &RepositoryFactory{
		db: db,
	}
}

// GetUserRepository UserRepository döndürür
func (f *RepositoryFactory) GetUserRepository() IUserRepository {
	f.mu.RLock()
	if f.userRepo != nil {
		defer f.mu.RUnlock()
		return f.userRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.userRepo == nil {
		f.userRepo = NewUserRepository(f.db)
	}
	return f.userRepo
}

// GetArticleRepository ArticleRepository döndürür
func (f *RepositoryFactory) GetArticleRepository() IArticleRepository {
	f.mu.RLock()
	if f.articleRepo != nil {
		defer f.mu.RUnlock()
		return f.articleRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.articleRepo == nil {
		f.articleRepo = NewArticleRepository(f.db)
	}
	return f.articleRepo
}

// GetCategoryRepository CategoryRepository döndürür
func (f *RepositoryFactory) GetCategoryRepository() ICategoryRepository {
	f.mu.RLock()
	if f.categoryRepo != nil {
		defer f.mu.RUnlock()
		return f.categoryRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.categoryRepo == nil {
		f.categoryRepo = NewCategoryRepository(f.db)
	}
	return f.categoryRepo
}

// GetTagRepository TagRepository döndürür
func (f *RepositoryFactory) GetTagRepository() ITagRepository {
	f.mu.RLock()
	if f.tagRepo != nil {
		defer f.mu.RUnlock()
		return f.tagRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.tagRepo == nil {
		f.tagRepo = NewTagRepository(f.db)
	}
	return f.tagRepo
}

// GetMediaRepository MediaRepository döndürür
func (f *RepositoryFactory) GetMediaRepository() IMediaRepository {
	f.mu.RLock()
	if f.mediaRepo != nil {
		defer f.mu.RUnlock()
		return f.mediaRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.mediaRepo == nil {
		f.mediaRepo = NewMediaRepository(f.db)
	}
	return f.mediaRepo
}

// SetUserRepository test için UserRepository'yi değiştirir
func (f *RepositoryFactory) SetUserRepository(repo IUserRepository) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.userRepo = repo
}

// SetArticleRepository test için ArticleRepository'yi değiştirir
func (f *RepositoryFactory) SetArticleRepository(repo IArticleRepository) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.articleRepo = repo
}

// SetCategoryRepository test için CategoryRepository'yi değiştirir
func (f *RepositoryFactory) SetCategoryRepository(repo ICategoryRepository) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.categoryRepo = repo
}

// SetTagRepository test için TagRepository'yi değiştirir
func (f *RepositoryFactory) SetTagRepository(repo ITagRepository) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.tagRepo = repo
}

// SetMediaRepository test için MediaRepository'yi değiştirir
func (f *RepositoryFactory) SetMediaRepository(repo IMediaRepository) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.mediaRepo = repo
}
