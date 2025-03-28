package domain

import (
	"time"

	"gorm.io/gorm"
)

// Comment yorum modelimiz
type Comment struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ArticleID  uint           `gorm:"not null" json:"article_id"`
	UserID     uint           `gorm:"not null" json:"user_id"`
	ParentID   *uint          `json:"parent_id,omitempty"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	IsApproved bool           `gorm:"default:false" json:"is_approved"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// İlişkiler
	User    *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Parent  *Comment   `gorm:"foreignKey:ParentID" json:"-"`
	Replies []*Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// CreateCommentRequest yorum oluşturma isteği
type CreateCommentRequest struct {
	ArticleID uint   `json:"article_id" validate:"required"`
	ParentID  *uint  `json:"parent_id,omitempty"`
	Content   string `json:"content" validate:"required"`
}

// UpdateCommentRequest yorum güncelleme isteği
type UpdateCommentRequest struct {
	Content    string `json:"content,omitempty"`
	IsApproved *bool  `json:"is_approved,omitempty"`
}
