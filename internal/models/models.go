package models

import (
	"time"

	"gorm.io/gorm"
)

// User kullanıcı modelimiz
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	FullName     string         `gorm:"size:100;not null" json:"full_name"`
	Role         string         `gorm:"size:20;not null;default:user" json:"role"` // admin, editor, user
	ProfileImage string         `gorm:"size:255" json:"profile_image,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// İlişkiler
	Articles []Article `gorm:"foreignKey:AuthorID" json:"-"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"-"`
}

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
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"-"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Articles []Article  `gorm:"foreignKey:CategoryID" json:"-"`
}

// Tag etiket modelimiz
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	Slug      string    `gorm:"size:50;uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time `json:"created_at"`

	// İlişkiler
	Articles []Article `gorm:"many2many:article_tags;" json:"-"`
}

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

	// İlişkiler
	Author   User      `json:"author,omitempty"`
	Category Category  `json:"category,omitempty"`
	Tags     []Tag     `gorm:"many2many:article_tags;" json:"tags,omitempty"`
	Comments []Comment `gorm:"foreignKey:ArticleID" json:"comments,omitempty"`
}

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
	Article Article   `json:"-"`
	User    User      `json:"user,omitempty"`
	Parent  *Comment  `gorm:"foreignKey:ParentID" json:"-"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// Media medya modelimiz
type Media struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Filename  string    `gorm:"size:255;not null" json:"filename"`
	Filepath  string    `gorm:"size:255;not null" json:"filepath"`
	Filetype  string    `gorm:"size:50;not null" json:"filetype"`
	Filesize  uint      `gorm:"not null" json:"filesize"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	// İlişkiler
	User User `json:"-"`
}

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
