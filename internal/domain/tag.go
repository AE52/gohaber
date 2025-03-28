package domain

import (
	"time"
)

// Tag etiket modelimiz
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	Slug      string    `gorm:"size:50;uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTagRequest etiket oluşturma isteği
type CreateTagRequest struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"required"`
}

// UpdateTagRequest etiket güncelleme isteği
type UpdateTagRequest struct {
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}
