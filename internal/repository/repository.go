package repository

import (
	"gorm.io/gorm"
)

// BaseRepository bütün repository'ler için ortak CRUD operasyonlarını tanımlar
type BaseRepository[T any, ID any] interface {
	Create(entity *T) error
	GetByID(id ID) (*T, error)
	Update(entity *T) error
	Delete(id ID) error
	List(offset, limit int, filters map[string]interface{}) ([]*T, int64, error)
	GetDB() *gorm.DB
}

// BaseRepositoryImpl BaseRepository'nin genel implementasyonu
type BaseRepositoryImpl[T any, ID any] struct {
	db        *gorm.DB
	entityPtr *T // GORM tarafından kullanılacak örnek entity
}

// NewBaseRepository yeni bir BaseRepository oluşturur
func NewBaseRepository[T any, ID any](db *gorm.DB, entityPtr *T) BaseRepository[T, ID] {
	return &BaseRepositoryImpl[T, ID]{
		db:        db,
		entityPtr: entityPtr,
	}
}

// Create yeni bir entity oluşturur
func (r *BaseRepositoryImpl[T, ID]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// GetByID ID'ye göre entity getirir
func (r *BaseRepositoryImpl[T, ID]) GetByID(id ID) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update entity'yi günceller
func (r *BaseRepositoryImpl[T, ID]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete entity'yi siler
func (r *BaseRepositoryImpl[T, ID]) Delete(id ID) error {
	return r.db.Delete(r.entityPtr, id).Error
}

// List entity'leri listeler
func (r *BaseRepositoryImpl[T, ID]) List(offset, limit int, filters map[string]interface{}) ([]*T, int64, error) {
	var entities []*T
	var count int64

	query := r.db.Model(r.entityPtr)

	// Filtreleri uygula
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, count, nil
}

// GetDB veritabanı bağlantısını döndürür
func (r *BaseRepositoryImpl[T, ID]) GetDB() *gorm.DB {
	return r.db
}
