package domain

import (
	"time"
)

// AdSpace reklam alanı modelimiz
type AdSpace struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"size:100;not null" json:"name"`
	Placement string     `gorm:"size:50;not null" json:"placement"` // header, sidebar, article-content, footer
	Content   string     `gorm:"type:text;not null" json:"content"` // HTML içeriği veya JS kodu
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// AdPlacement reklam yerleşim sabitleri
const (
	AdPlacementHeader        = "header"
	AdPlacementSidebar       = "sidebar"
	AdPlacementArticleTop    = "article-top"
	AdPlacementArticleBottom = "article-bottom"
	AdPlacementArticleMiddle = "article-middle"
	AdPlacementFooter        = "footer"
)

// CreateAdSpaceRequest reklam alanı oluşturma isteği
type CreateAdSpaceRequest struct {
	Name      string `json:"name" validate:"required"`
	Placement string `json:"placement" validate:"required,oneof=header sidebar article-top article-bottom article-middle footer"`
	Content   string `json:"content" validate:"required"`
	IsActive  bool   `json:"is_active"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

// UpdateAdSpaceRequest reklam alanı güncelleme isteği
type UpdateAdSpaceRequest struct {
	Name      string `json:"name,omitempty"`
	Placement string `json:"placement,omitempty" validate:"omitempty,oneof=header sidebar article-top article-bottom article-middle footer"`
	Content   string `json:"content,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}
