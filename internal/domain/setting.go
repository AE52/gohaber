package domain

import (
	"time"
)

// Setting ayar modelimiz
type Setting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	Group     string    `gorm:"size:50;not null;index" json:"group"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SettingGroup ayar grubu sabitleri
const (
	SettingGroupGeneral     = "general"
	SettingGroupAppearance  = "appearance"
	SettingGroupIntegration = "integration"
	SettingGroupEmail       = "email"
	SettingGroupSocial      = "social"
	SettingGroupSEO         = "seo"
	SettingGroupCache       = "cache"
	SettingGroupBackup      = "backup"
	SettingGroupAdmin       = "admin"
)

// SettingType ayar tipi sabitleri
const (
	SettingTypeString  = "string"
	SettingTypeNumber  = "number"
	SettingTypeBoolean = "boolean"
	SettingTypeJSON    = "json"
)

// AllSettings, tüm ayarları içeren yapıdır
type AllSettings struct {
	// Genel Ayarlar
	SiteName            string `json:"site_name"`
	SiteDescription     string `json:"site_description"`
	SiteURL             string `json:"site_url"`
	AdminEmail          string `json:"admin_email"`
	Timezone            string `json:"timezone"`
	DateFormat          string `json:"date_format"`
	EnableComments      bool   `json:"enable_comments"`
	AutoApproveComments bool   `json:"auto_approve_comments"`
	DefaultLanguage     string `json:"default_language"`

	// Görünüm Ayarları
	Theme          string `json:"theme"`
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	Logo           string `json:"logo"`
	Favicon        string `json:"favicon"`
	HomepageLayout string `json:"homepage_layout"`
	CustomCSS      string `json:"custom_css"`

	// Entegrasyon Ayarları
	GoogleAnalyticsID  string `json:"google_analytics_id"`
	GoogleTagManagerID string `json:"google_tag_manager_id"`
	RecaptchaSiteKey   string `json:"recaptcha_site_key"`
	RecaptchaSecretKey string `json:"recaptcha_secret_key"`
	FacebookPixelID    string `json:"facebook_pixel_id"`
	DisqusShortname    string `json:"disqus_shortname"`
	CustomHeader       string `json:"custom_header"`
	CustomFooter       string `json:"custom_footer"`

	// Diğer ayar grupları da ihtiyaca göre eklenebilir
}

// SettingRequest ayar isteği
type SettingRequest struct {
	Group string      `json:"group" validate:"required"`
	Key   string      `json:"key" validate:"required"`
	Value interface{} `json:"value"`
	Type  string      `json:"type" validate:"required,oneof=string number boolean json"`
}

// GeneralSettings genel site ayarları
type GeneralSettings struct {
	SiteName        string `json:"site_name"`
	SiteDescription string `json:"site_description"`
	SiteURL         string `json:"site_url"`
	SiteLogo        string `json:"site_logo"`
	SiteFavicon     string `json:"site_favicon"`
	ContactEmail    string `json:"contact_email"`
	ContactPhone    string `json:"contact_phone"`
	FooterText      string `json:"footer_text"`
	MetaKeywords    string `json:"meta_keywords"`
	MetaDescription string `json:"meta_description"`
}

// SocialSettings sosyal medya ayarları
type SocialSettings struct {
	Facebook  string `json:"facebook"`
	Twitter   string `json:"twitter"`
	Instagram string `json:"instagram"`
	LinkedIn  string `json:"linkedin"`
	YouTube   string `json:"youtube"`
}
