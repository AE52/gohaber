package domain

import (
	"time"
)

// Media medya modelimiz
type Media struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Filename    string    `gorm:"size:255;not null" json:"filename"`
	ObjectName  string    `gorm:"size:255;not null" json:"object_name"` // MinIO nesne adı
	ContentType string    `gorm:"size:100;not null" json:"content_type"`
	Filesize    uint      `gorm:"not null" json:"filesize"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// İlişkiler
	User *User `json:"-" gorm:"foreignKey:UserID"`
}

// MediaResponse medya yanıtı
type MediaResponse struct {
	ID          uint      `json:"id"`
	Filename    string    `json:"filename"`
	ObjectName  string    `json:"object_name"`
	ContentType string    `json:"content_type"`
	Filesize    uint      `json:"filesize"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
}
