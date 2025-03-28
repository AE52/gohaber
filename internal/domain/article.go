package domain

import (
	"time"

	"gorm.io/gorm"
)

// Article haber/makale modelimiz
type Article struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Slug          string         `gorm:"size:200;uniqueIndex;not null" json:"slug"`
	Content       string         `gorm:"type:text;not null" json:"content"`
	Summary       string         `gorm:"type:text;not null" json:"summary"`
	FeaturedImage string         `gorm:"size:255" json:"featured_image,omitempty"`
	AuthorID      uint           `gorm:"not null" json:"author_id"`
	CategoryID    uint           `gorm:"not null" json:"category_id"`
	Status        string         `gorm:"size:20;not null;default:draft" json:"status"` // published, draft, pending
	ViewCount     uint           `gorm:"default:0" json:"view_count"`
	IsFeatured    bool           `gorm:"default:false" json:"is_featured"`
	PublishedAt   *time.Time     `json:"published_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// İlişkiler - JSON dönüşümünde çözünürlük için
	Author   *User      `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Category *Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Tags     []*Tag     `gorm:"many2many:article_tags;" json:"tags,omitempty"`
	Comments []*Comment `json:"comments,omitempty" gorm:"foreignKey:ArticleID"`
}

// ArticleStatus makale durum sabitleri
const (
	ArticleStatusPublished = "published"
	ArticleStatusDraft     = "draft"
	ArticleStatusPending   = "pending"
)

// CreateArticleRequest makale oluşturma isteği
type CreateArticleRequest struct {
	Title         string `json:"title" validate:"required"`
	Content       string `json:"content" validate:"required"`
	Summary       string `json:"summary" validate:"required"`
	FeaturedImage string `json:"featured_image,omitempty"`
	CategoryID    uint   `json:"category_id" validate:"required"`
	Status        string `json:"status" validate:"oneof=published draft pending"`
	IsFeatured    bool   `json:"is_featured"`
	TagIDs        []uint `json:"tag_ids,omitempty"`
	PublishedAt   string `json:"published_at,omitempty"`
}

// UpdateArticleRequest makale güncelleme isteği
type UpdateArticleRequest struct {
	Title         string `json:"title,omitempty"`
	Content       string `json:"content,omitempty"`
	Summary       string `json:"summary,omitempty"`
	FeaturedImage string `json:"featured_image,omitempty"`
	CategoryID    uint   `json:"category_id,omitempty"`
	Status        string `json:"status,omitempty" validate:"omitempty,oneof=published draft pending"`
	IsFeatured    *bool  `json:"is_featured,omitempty"`
	TagIDs        []uint `json:"tag_ids,omitempty"`
	PublishedAt   string `json:"published_at,omitempty"`
}
