package repository

import (
	"log"

	"github.com/username/haber/internal/config"
	"github.com/username/haber/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database veritabanı işlemlerini kapsüller
type Database struct {
	DB *gorm.DB
}

// NewDatabase yeni veritabanı bağlantısı oluşturur
func NewDatabase(config config.IDatabaseConfig) *Database {
	dsn := config.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
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

// AutoMigrate veritabanı şemasını otomatik günceller
func (d *Database) AutoMigrate() error {
	return d.DB.AutoMigrate(
		&domain.User{},
		&domain.Category{},
		&domain.Tag{},
		&domain.Article{},
		&domain.Comment{},
		&domain.Media{},
		&domain.Setting{},
		&domain.AdSpace{},
	)
}

// WithTransaction transaction başlatır ve işler
func (d *Database) WithTransaction(fn func(tx *gorm.DB) error) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
