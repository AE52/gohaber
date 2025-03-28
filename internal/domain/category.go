package domain

import (
	"time"
)

// Category kategori modelimiz
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Slug        string    `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Description string    `json:"description,omitempty"`
	ParentID    *uint     `json:"parent_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// İlişkiler
	Parent   *Category   `gorm:"foreignKey:ParentID" json:"-"`
	Children []*Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Articles []*Article  `gorm:"foreignKey:CategoryID" json:"-"`
}

// CreateCategoryRequest kategori oluşturma isteği
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"`
}

// UpdateCategoryRequest kategori güncelleme isteği
type UpdateCategoryRequest struct {
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"`
}
