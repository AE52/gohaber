package repository

import (
	"log"

	"github.com/username/haber/internal/config"
	"github.com/username/haber/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database veritabanı bağlantısını yönetir
type Database struct {
	DB *gorm.DB
}

// NewDatabase yeni bir Database örneği oluşturur
func NewDatabase(cfg *config.DatabaseConfig) *Database {
	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Veritabanı bağlantısı kurulamadı: %v", err)
	}

	return &Database{DB: db}
}

// Close veritabanı bağlantısını kapatır
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate veritabanı tablolarını otomatik olarak oluşturur
func (d *Database) AutoMigrate() error {
	return d.DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Tag{},
		&models.Article{},
		&models.Comment{},
		&models.Media{},
		&models.AdSpace{},
	)
}

// UserRepository kullanıcı işlemleri için repository
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository yeni bir UserRepository örneği oluşturur
func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db.DB}
}

// ArticleRepository makale işlemleri için repository
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository yeni bir ArticleRepository örneği oluşturur
func NewArticleRepository(db *Database) *ArticleRepository {
	return &ArticleRepository{db: db.DB}
}

// CategoryRepository kategori işlemleri için repository
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository yeni bir CategoryRepository örneği oluşturur
func NewCategoryRepository(db *Database) *CategoryRepository {
	return &CategoryRepository{db: db.DB}
}

// TagRepository etiket işlemleri için repository
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository yeni bir TagRepository örneği oluşturur
func NewTagRepository(db *Database) *TagRepository {
	return &TagRepository{db: db.DB}
}

// CommentRepository yorum işlemleri için repository
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository yeni bir CommentRepository örneği oluşturur
func NewCommentRepository(db *Database) *CommentRepository {
	return &CommentRepository{db: db.DB}
}

// MediaRepository medya işlemleri için repository
type MediaRepository struct {
	db *gorm.DB
}

// NewMediaRepository yeni bir MediaRepository örneği oluşturur
func NewMediaRepository(db *Database) *MediaRepository {
	return &MediaRepository{db: db.DB}
}

// AdSpaceRepository reklam alanları için repository
type AdSpaceRepository struct {
	db *gorm.DB
}

// NewAdSpaceRepository yeni bir AdSpaceRepository örneği oluşturur
func NewAdSpaceRepository(db *Database) *AdSpaceRepository {
	return &AdSpaceRepository{db: db.DB}
}
