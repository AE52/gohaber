package domain

import "time"

// MetaData meta veri yapısı
type MetaData map[string]interface{}

// APIResponse genel API yanıt yapısı
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    MetaData    `json:"meta,omitempty"`
}

// APIError API hata yapısı
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse swagger dokümantasyonu için standart hata yanıtı
type ErrorResponse struct {
	Error   string `json:"error" example:"Hata açıklaması"`
	Status  int    `json:"status" example:"400"`
	Message string `json:"message,omitempty" example:"Detaylı hata mesajı"`
}

// PaginatedResponse sayfalanmış yanıt yapısı
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
	PerPage    int         `json:"per_page"`
	Page       int         `json:"page"`
}

// AuthResponse kimlik doğrulama yanıtı
type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// TokenResponse token yanıtı
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// UploadInfo yükleme bilgisi yapısı
type UploadInfo struct {
	ID         uint      `json:"id"`
	Filename   string    `json:"filename"`
	Path       string    `json:"path"`
	URL        string    `json:"url"`
	Size       int64     `json:"size"`
	Type       string    `json:"type"`
	MIME       string    `json:"mime"`
	UploadType string    `json:"upload_type"`
	UserID     uint      `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MessageResponse basit mesaj yanıtı
type MessageResponse struct {
	Message string `json:"message" example:"İşlem başarıyla tamamlandı"`
}
