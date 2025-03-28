package service

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/username/haber/internal/domain"
	"gorm.io/gorm"
)

// Yardımcı fonksiyonlar
func getOrDefault(m map[string]string, key, defaultValue string) string {
	if value, exists := m[key]; exists && value != "" {
		return value
	}
	return defaultValue
}

func getBoolOrDefault(m map[string]string, key string, defaultValue bool) bool {
	if value, exists := m[key]; exists {
		if value == "1" || value == "true" {
			return true
		} else if value == "0" || value == "false" {
			return false
		}
	}
	return defaultValue
}

func getIntOrDefault(m map[string]string, key string, defaultValue int) int {
	if value, exists := m[key]; exists && value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// ISettingsService, ayarlar ile ilgili veritabanı işlemleri için arayüz
type ISettingsService interface {
	GetSetting(key string) (string, error)
	SetSetting(key, value, group string) error
	GetSettingsByGroup(group string) (map[string]string, error)
	GetAllSettings() (*domain.AllSettings, error)
	SaveGeneralSettings(settings map[string]interface{}) error
	SaveAppearanceSettings(settings map[string]interface{}) error
	SaveIntegrationSettings(settings map[string]interface{}) error
	SaveEmailSettings(settings map[string]interface{}) error
	SaveSocialSettings(settings map[string]interface{}) error
	SaveSEOSettings(settings map[string]interface{}) error
	SaveCacheSettings(settings map[string]interface{}) error
	SaveBackupSettings(settings map[string]interface{}) error
}

// SettingsService, ayarlar servisi implementasyonu
type SettingsService struct {
	DB *gorm.DB
}

// NewSettingsService yeni bir SettingsService oluşturur
func NewSettingsService(db *gorm.DB) ISettingsService {
	return &SettingsService{
		DB: db,
	}
}

// GetSetting, belirli bir ayarı getirir
func (s *SettingsService) GetSetting(key string) (string, error) {
	var setting domain.Setting
	err := s.DB.Where("key = ?", key).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return setting.Value, nil
}

// SetSetting, belirli bir ayarı kaydeder
func (s *SettingsService) SetSetting(key, value, group string) error {
	now := time.Now()

	var setting domain.Setting
	err := s.DB.Where("key = ?", key).First(&setting).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Kayıt yok, yeni oluştur
			setting = domain.Setting{
				Key:       key,
				Value:     value,
				Group:     group,
				CreatedAt: now,
				UpdatedAt: now,
			}
			return s.DB.Create(&setting).Error
		}
		return err
	}

	// Mevcut kaydı güncelle
	setting.Value = value
	setting.Group = group
	setting.UpdatedAt = now
	return s.DB.Save(&setting).Error
}

// GetSettingsByGroup, belirli bir gruptaki tüm ayarları getirir
func (s *SettingsService) GetSettingsByGroup(group string) (map[string]string, error) {
	var settings []domain.Setting
	err := s.DB.Where("group = ?", group).Find(&settings).Error
	if err != nil {
		return nil, err
	}

	settingsMap := make(map[string]string)
	for _, setting := range settings {
		settingsMap[setting.Key] = setting.Value
	}

	return settingsMap, nil
}

// GetAllSettings, tüm ayarları getirir
func (s *SettingsService) GetAllSettings() (*domain.AllSettings, error) {
	var settings []domain.Setting
	err := s.DB.Find(&settings).Error
	if err != nil {
		return nil, err
	}

	settingsMap := make(map[string]string)
	for _, setting := range settings {
		settingsMap[setting.Key] = setting.Value
	}

	// Varsayılan ayarları içeren AllSettings nesnesi oluştur
	result := &domain.AllSettings{
		SiteName:            settingsMap["site_name"],
		SiteDescription:     settingsMap["site_description"],
		SiteURL:             settingsMap["site_url"],
		AdminEmail:          settingsMap["admin_email"],
		Timezone:            getOrDefault(settingsMap, "timezone", "Europe/Istanbul"),
		DateFormat:          getOrDefault(settingsMap, "date_format", "DD.MM.YYYY"),
		EnableComments:      getBoolOrDefault(settingsMap, "enable_comments", true),
		AutoApproveComments: getBoolOrDefault(settingsMap, "auto_approve_comments", false),
		DefaultLanguage:     getOrDefault(settingsMap, "default_language", "tr"),
	}

	// Tüm ayar alanlarını doldur - burada eksik kalan ayarları doldururuz
	// Bu kısımda domain.AllSettings yapısının tüm alanlarını doldurmak gerekir

	return result, nil
}

// SaveGeneralSettings, genel ayarları kaydeder
func (s *SettingsService) SaveGeneralSettings(settings map[string]interface{}) error {
	return s.saveSettingsToGorm(settings, "general")
}

// SaveAppearanceSettings, görünüm ayarlarını kaydeder
func (s *SettingsService) SaveAppearanceSettings(settings map[string]interface{}) error {
	return s.saveSettingsToGorm(settings, "appearance")
}

// SaveIntegrationSettings, entegrasyon ayarlarını kaydeder
func (s *SettingsService) SaveIntegrationSettings(settings map[string]interface{}) error {
	return s.saveSettingsToGorm(settings, "integration")
}

// SaveEmailSettings, e-posta ayarlarını kaydeder
func (s *SettingsService) SaveEmailSettings(settings map[string]interface{}) error {
	return s.saveSettingsToGorm(settings, "email")
}

// SaveSocialSettings, sosyal medya ayarlarını kaydeder
func (s *SettingsService) SaveSocialSettings(settings map[string]interface{}) error {
	return s.saveSettingsToGorm(settings, "social")
}

// SaveSEOSettings, SEO ayarlarını kaydeder
func (s *SettingsService) SaveSEOSettings(settings map[string]interface{}) error {
	return s.saveSettingsToGorm(settings, "seo")
}

// SaveCacheSettings, önbellek ayarlarını kaydeder
func (s *SettingsService) SaveCacheSettings(settings map[string]interface{}) error {
	return s.saveSettingsToGorm(settings, "cache")
}

// SaveBackupSettings, yedekleme ayarlarını kaydeder
func (s *SettingsService) SaveBackupSettings(settings map[string]interface{}) error {
	return s.saveSettingsToGorm(settings, "backup")
}

// GORM için yardımcı fonksiyon
func (s *SettingsService) saveSettingsToGorm(settings map[string]interface{}, group string) error {
	for key, value := range settings {
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case bool:
			if v {
				strValue = "1"
			} else {
				strValue = "0"
			}
		default:
			bytes, err := json.Marshal(v)
			if err != nil {
				return err
			}
			strValue = string(bytes)
		}

		if err := s.SetSetting(key, strValue, group); err != nil {
			return err
		}
	}
	return nil
}
