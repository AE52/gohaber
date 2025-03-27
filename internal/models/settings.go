package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Settings, site ayarlarını temsil eden ana yapıdır
type Settings struct {
	ID        int64     `json:"id" db:"id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	Group     string    `json:"group" db:"group"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SettingsGroup, ayarların gruplarını temsil eder
type SettingsGroup string

const (
	GeneralSettings     SettingsGroup = "general"
	AppearanceSettings  SettingsGroup = "appearance"
	IntegrationSettings SettingsGroup = "integration"
	EmailSettings       SettingsGroup = "email"
	SocialSettings      SettingsGroup = "social"
	SEOSettings         SettingsGroup = "seo"
	CacheSettings       SettingsGroup = "cache"
	BackupSettings      SettingsGroup = "backup"
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

	// E-posta Ayarları
	EmailProvider string `json:"email_provider"`
	SMTPHost      string `json:"smtp_host"`
	SMTPPort      int    `json:"smtp_port"`
	SMTPSecurity  string `json:"smtp_security"`
	SMTPUsername  string `json:"smtp_username"`
	SMTPPassword  string `json:"smtp_password"`
	EmailAPIKey   string `json:"email_api_key"`
	FromEmail     string `json:"from_email"`
	FromName      string `json:"from_name"`

	// Sosyal Medya Ayarları
	FacebookURL         string `json:"facebook_url"`
	TwitterURL          string `json:"twitter_url"`
	InstagramURL        string `json:"instagram_url"`
	YouTubeURL          string `json:"youtube_url"`
	LinkedInURL         string `json:"linkedin_url"`
	EnableSocialSharing bool   `json:"enable_social_sharing"`
	ShareFacebook       bool   `json:"share_facebook"`
	ShareTwitter        bool   `json:"share_twitter"`
	ShareWhatsapp       bool   `json:"share_whatsapp"`
	ShareLinkedin       bool   `json:"share_linkedin"`
	ShareTelegram       bool   `json:"share_telegram"`
	ShareEmail          bool   `json:"share_email"`
	EnableSocialLogin   bool   `json:"enable_social_login"`
	EnableFacebookLogin bool   `json:"enable_facebook_login"`
	FacebookAppID       string `json:"facebook_app_id"`
	FacebookAppSecret   string `json:"facebook_app_secret"`
	EnableGoogleLogin   bool   `json:"enable_google_login"`
	GoogleClientID      string `json:"google_client_id"`
	GoogleClientSecret  string `json:"google_client_secret"`

	// SEO Ayarları
	HomePageTitle          string `json:"home_page_title"`
	ArticleTitle           string `json:"article_title"`
	CategoryTitle          string `json:"category_title"`
	TagTitle               string `json:"tag_title"`
	AuthorTitle            string `json:"author_title"`
	MetaDescription        string `json:"meta_description"`
	ArticleMetaDescription string `json:"article_meta_description"`
	EnableSitemap          bool   `json:"enable_sitemap"`
	SitemapChangeFreq      string `json:"sitemap_change_freq"`
	SitemapPriority        string `json:"sitemap_priority"`
	RobotsTxt              string `json:"robots_txt"`
	EnableOpenGraph        bool   `json:"enable_open_graph"`
	EnableTwitterCards     bool   `json:"enable_twitter_cards"`
	TwitterUsername        string `json:"twitter_username"`
	DefaultSocialImage     string `json:"default_social_image"`
	EnableStructuredData   bool   `json:"enable_structured_data"`
	OrganizationType       string `json:"organization_type"`
	OrganizationLogo       string `json:"organization_logo"`
	ArticleType            string `json:"article_type"`
	EnableCanonicalUrls    bool   `json:"enable_canonical_urls"`
	NoindexCategories      bool   `json:"noindex_categories"`
	NoindexTags            bool   `json:"noindex_tags"`
	NoindexAuthors         bool   `json:"noindex_authors"`
	NoindexArchives        bool   `json:"noindex_archives"`

	// Önbellek Ayarları
	EnablePageCache         bool   `json:"enable_page_cache"`
	PageCacheDuration       int    `json:"page_cache_duration"`
	ExcludedCachePaths      string `json:"excluded_cache_paths"`
	EnableQueryCache        bool   `json:"enable_query_cache"`
	QueryCacheDuration      int    `json:"query_cache_duration"`
	EnableBrowserCache      bool   `json:"enable_browser_cache"`
	BrowserCacheDuration    int    `json:"browser_cache_duration"`
	EnableCDN               bool   `json:"enable_cdn"`
	CDNURL                  string `json:"cdn_url"`
	CDNImages               bool   `json:"cdn_images"`
	CDNCSS                  bool   `json:"cdn_css"`
	CDNJS                   bool   `json:"cdn_js"`
	CDNFonts                bool   `json:"cdn_fonts"`
	EnableImageOptimization bool   `json:"enable_image_optimization"`
	ImageQuality            int    `json:"image_quality"`
	EnableWebP              bool   `json:"enable_webp"`
	EnableLazyLoading       bool   `json:"enable_lazy_loading"`
	EnableJSMinify          bool   `json:"enable_js_minify"`
	EnableCSSMinify         bool   `json:"enable_css_minify"`
	EnableHTMLMinify        bool   `json:"enable_html_minify"`
	CombineJSFiles          bool   `json:"combine_js_files"`
	CombineCSSFiles         bool   `json:"combine_css_files"`

	// Yedekleme Ayarları
	EnableAutoBackup        bool   `json:"enable_auto_backup"`
	BackupFrequency         string `json:"backup_frequency"`
	BackupDay               string `json:"backup_day"`
	BackupDate              string `json:"backup_date"`
	BackupTime              string `json:"backup_time"`
	BackupDatabase          bool   `json:"backup_database"`
	BackupFiles             bool   `json:"backup_files"`
	BackupSettings          bool   `json:"backup_settings"`
	BackupStorage           string `json:"backup_storage"`
	LocalBackupPath         string `json:"local_backup_path"`
	CloudAPIKey             string `json:"cloud_api_key"`
	CloudSecret             string `json:"cloud_secret"`
	CloudFolder             string `json:"cloud_folder"`
	BackupRetention         int    `json:"backup_retention"`
	CompressBackups         bool   `json:"compress_backups"`
	EncryptBackups          bool   `json:"encrypt_backups"`
	BackupNotifications     bool   `json:"backup_notifications"`
	BackupNotificationEmail string `json:"backup_notification_email"`
}

// Backup, yedekleme kayıtlarını temsil eder
type Backup struct {
	ID        int64     `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"`
	Size      string    `json:"size" db:"size"`
	Path      string    `json:"path" db:"path"`
	Storage   string    `json:"storage" db:"storage"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// SettingsService, ayarlar ile ilgili veritabanı işlemleri için arayüz
type SettingsService interface {
	GetSetting(key string) (string, error)
	SetSetting(key, value, group string) error
	GetSettingsByGroup(group string) (map[string]string, error)
	GetAllSettings() (*AllSettings, error)
	SaveGeneralSettings(settings map[string]interface{}) error
	SaveAppearanceSettings(settings map[string]interface{}) error
	SaveIntegrationSettings(settings map[string]interface{}) error
	SaveEmailSettings(settings map[string]interface{}) error
	SaveSocialSettings(settings map[string]interface{}) error
	SaveSEOSettings(settings map[string]interface{}) error
	SaveCacheSettings(settings map[string]interface{}) error
	SaveBackupSettings(settings map[string]interface{}) error
}

// SqlSettingsService, SettingsService arayüzünü SQL implementasyonu ile gerçekleştirir
type SqlSettingsService struct {
	DB *sql.DB
}

// GetSetting, belirli bir ayarı getirir
func (s *SqlSettingsService) GetSetting(key string) (string, error) {
	var setting Settings
	query := `SELECT * FROM settings WHERE key = ?`
	err := s.DB.QueryRow(query, key).Scan(
		&setting.ID,
		&setting.Key,
		&setting.Value,
		&setting.Group,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return setting.Value, nil
}

// SetSetting, belirli bir ayarı kaydeder
func (s *SqlSettingsService) SetSetting(key, value, group string) error {
	now := time.Now()
	query := `INSERT INTO settings (key, value, group, created_at, updated_at) 
	VALUES (?, ?, ?, ?, ?) 
	ON DUPLICATE KEY UPDATE value = ?, group = ?, updated_at = ?`

	_, err := s.DB.Exec(query, key, value, group, now, now, value, group, now)
	return err
}

// GetSettingsByGroup, belirli bir gruptaki tüm ayarları getirir
func (s *SqlSettingsService) GetSettingsByGroup(group string) (map[string]string, error) {
	query := `SELECT key, value FROM settings WHERE group = ?`
	rows, err := s.DB.Query(query, group)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		settings[key] = value
	}

	return settings, nil
}

// GetAllSettings, tüm ayarları getirir
func (s *SqlSettingsService) GetAllSettings() (*AllSettings, error) {
	query := `SELECT key, value FROM settings`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settingsMap := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		settingsMap[key] = value
	}

	// Varsayılan ayarları içeren AllSettings nesnesi oluştur
	settings := &AllSettings{
		SiteName:            settingsMap["site_name"],
		SiteDescription:     settingsMap["site_description"],
		SiteURL:             settingsMap["site_url"],
		AdminEmail:          settingsMap["admin_email"],
		Timezone:            getOrDefault(settingsMap, "timezone", "Europe/Istanbul"),
		DateFormat:          getOrDefault(settingsMap, "date_format", "DD.MM.YYYY"),
		EnableComments:      getBoolOrDefault(settingsMap, "enable_comments", true),
		AutoApproveComments: getBoolOrDefault(settingsMap, "auto_approve_comments", false),
		DefaultLanguage:     getOrDefault(settingsMap, "default_language", "tr"),

		// Görünüm Ayarları
		Theme:          getOrDefault(settingsMap, "theme", "default"),
		PrimaryColor:   getOrDefault(settingsMap, "primary_color", "#007bff"),
		SecondaryColor: getOrDefault(settingsMap, "secondary_color", "#6c757d"),
		Logo:           settingsMap["logo"],
		Favicon:        settingsMap["favicon"],
		HomepageLayout: getOrDefault(settingsMap, "homepage_layout", "default"),
		CustomCSS:      settingsMap["custom_css"],

		// Diğer tüm ayarlar benzer şekilde doldurulacak
		// ...
	}

	// Tüm ayar alanlarını doldur
	fillSettings(settings, settingsMap)

	return settings, nil
}

// Ayarları kaydetmek için yardımcı fonksiyonlar
func (s *SqlSettingsService) SaveGeneralSettings(settings map[string]interface{}) error {
	// Bu map'i işleyerek veritabanına kaydet
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

		if err := s.SetSetting(key, strValue, string(GeneralSettings)); err != nil {
			return err
		}
	}
	return nil
}

// Diğer kaydetme fonksiyonları da benzer şekilde
func (s *SqlSettingsService) SaveAppearanceSettings(settings map[string]interface{}) error {
	return saveSettingsGroup(s, settings, string(AppearanceSettings))
}

func (s *SqlSettingsService) SaveIntegrationSettings(settings map[string]interface{}) error {
	return saveSettingsGroup(s, settings, string(IntegrationSettings))
}

func (s *SqlSettingsService) SaveEmailSettings(settings map[string]interface{}) error {
	return saveSettingsGroup(s, settings, string(EmailSettings))
}

func (s *SqlSettingsService) SaveSocialSettings(settings map[string]interface{}) error {
	return saveSettingsGroup(s, settings, string(SocialSettings))
}

func (s *SqlSettingsService) SaveSEOSettings(settings map[string]interface{}) error {
	return saveSettingsGroup(s, settings, string(SEOSettings))
}

func (s *SqlSettingsService) SaveCacheSettings(settings map[string]interface{}) error {
	return saveSettingsGroup(s, settings, string(CacheSettings))
}

func (s *SqlSettingsService) SaveBackupSettings(settings map[string]interface{}) error {
	return saveSettingsGroup(s, settings, string(BackupSettings))
}

// Yardımcı fonksiyonlar
func saveSettingsGroup(s *SqlSettingsService, settings map[string]interface{}, group string) error {
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

func getOrDefault(m map[string]string, key, defaultValue string) string {
	if val, exists := m[key]; exists && val != "" {
		return val
	}
	return defaultValue
}

func getBoolOrDefault(m map[string]string, key string, defaultValue bool) bool {
	if val, exists := m[key]; exists {
		return val == "1" || val == "true"
	}
	return defaultValue
}

func fillSettings(settings *AllSettings, settingsMap map[string]string) {
	// Bu fonksiyon settingsMap'ten AllSettings nesnesinin tüm alanlarını doldurur
	// Yukarıda bazı alanlar manuel olarak dolduruldu, geri kalanı buradan doldurulabilir

	// Entegrasyon Ayarları
	settings.GoogleAnalyticsID = settingsMap["google_analytics_id"]
	settings.GoogleTagManagerID = settingsMap["google_tag_manager_id"]
	settings.RecaptchaSiteKey = settingsMap["recaptcha_site_key"]
	settings.RecaptchaSecretKey = settingsMap["recaptcha_secret_key"]
	settings.FacebookPixelID = settingsMap["facebook_pixel_id"]
	settings.DisqusShortname = settingsMap["disqus_shortname"]
	settings.CustomHeader = settingsMap["custom_header"]
	settings.CustomFooter = settingsMap["custom_footer"]

	// E-posta Ayarları
	settings.EmailProvider = getOrDefault(settingsMap, "email_provider", "smtp")
	settings.SMTPHost = settingsMap["smtp_host"]
	settings.SMTPPort = getIntOrDefault(settingsMap, "smtp_port", 587)
	settings.SMTPSecurity = getOrDefault(settingsMap, "smtp_security", "tls")
	settings.SMTPUsername = settingsMap["smtp_username"]
	settings.SMTPPassword = settingsMap["smtp_password"]
	settings.EmailAPIKey = settingsMap["email_api_key"]
	settings.FromEmail = settingsMap["from_email"]
	settings.FromName = settingsMap["from_name"]

	// Sosyal Medya Ayarları
	settings.FacebookURL = settingsMap["facebook_url"]
	settings.TwitterURL = settingsMap["twitter_url"]
	settings.InstagramURL = settingsMap["instagram_url"]
	settings.YouTubeURL = settingsMap["youtube_url"]
	settings.LinkedInURL = settingsMap["linkedin_url"]
	settings.EnableSocialSharing = getBoolOrDefault(settingsMap, "enable_social_sharing", true)
	settings.ShareFacebook = getBoolOrDefault(settingsMap, "share_facebook", true)
	settings.ShareTwitter = getBoolOrDefault(settingsMap, "share_twitter", true)
	settings.ShareWhatsapp = getBoolOrDefault(settingsMap, "share_whatsapp", true)
	settings.ShareLinkedin = getBoolOrDefault(settingsMap, "share_linkedin", true)
	settings.ShareTelegram = getBoolOrDefault(settingsMap, "share_telegram", false)
	settings.ShareEmail = getBoolOrDefault(settingsMap, "share_email", true)
	settings.EnableSocialLogin = getBoolOrDefault(settingsMap, "enable_social_login", false)
	settings.EnableFacebookLogin = getBoolOrDefault(settingsMap, "enable_facebook_login", false)
	settings.FacebookAppID = settingsMap["facebook_app_id"]
	settings.FacebookAppSecret = settingsMap["facebook_app_secret"]
	settings.EnableGoogleLogin = getBoolOrDefault(settingsMap, "enable_google_login", false)
	settings.GoogleClientID = settingsMap["google_client_id"]
	settings.GoogleClientSecret = settingsMap["google_client_secret"]

	// SEO Ayarları
	settings.HomePageTitle = getOrDefault(settingsMap, "home_page_title", "{site_name} - {site_description}")
	settings.ArticleTitle = getOrDefault(settingsMap, "article_title", "{page_title} - {site_name}")
	settings.CategoryTitle = getOrDefault(settingsMap, "category_title", "{category} - {site_name}")
	settings.TagTitle = getOrDefault(settingsMap, "tag_title", "{tag} - {site_name}")
	settings.AuthorTitle = getOrDefault(settingsMap, "author_title", "{author} - {site_name}")
	settings.MetaDescription = settingsMap["meta_description"]
	settings.ArticleMetaDescription = settingsMap["article_meta_description"]
	settings.EnableSitemap = getBoolOrDefault(settingsMap, "enable_sitemap", true)
	settings.SitemapChangeFreq = getOrDefault(settingsMap, "sitemap_change_freq", "daily")
	settings.SitemapPriority = getOrDefault(settingsMap, "sitemap_priority", "0.7")
	settings.RobotsTxt = settingsMap["robots_txt"]
	settings.EnableOpenGraph = getBoolOrDefault(settingsMap, "enable_open_graph", true)
	settings.EnableTwitterCards = getBoolOrDefault(settingsMap, "enable_twitter_cards", true)
	settings.TwitterUsername = settingsMap["twitter_username"]
	settings.DefaultSocialImage = settingsMap["default_social_image"]
	settings.EnableStructuredData = getBoolOrDefault(settingsMap, "enable_structured_data", true)
	settings.OrganizationType = getOrDefault(settingsMap, "organization_type", "NewsMediaOrganization")
	settings.OrganizationLogo = settingsMap["organization_logo"]
	settings.ArticleType = getOrDefault(settingsMap, "article_type", "NewsArticle")
	settings.EnableCanonicalUrls = getBoolOrDefault(settingsMap, "enable_canonical_urls", true)
	settings.NoindexCategories = getBoolOrDefault(settingsMap, "noindex_categories", false)
	settings.NoindexTags = getBoolOrDefault(settingsMap, "noindex_tags", false)
	settings.NoindexAuthors = getBoolOrDefault(settingsMap, "noindex_authors", false)
	settings.NoindexArchives = getBoolOrDefault(settingsMap, "noindex_archives", false)

	// Önbellek Ayarları
	settings.EnablePageCache = getBoolOrDefault(settingsMap, "enable_page_cache", true)
	settings.PageCacheDuration = getIntOrDefault(settingsMap, "page_cache_duration", 60)
	settings.ExcludedCachePaths = settingsMap["excluded_cache_paths"]
	settings.EnableQueryCache = getBoolOrDefault(settingsMap, "enable_query_cache", true)
	settings.QueryCacheDuration = getIntOrDefault(settingsMap, "query_cache_duration", 30)
	settings.EnableBrowserCache = getBoolOrDefault(settingsMap, "enable_browser_cache", true)
	settings.BrowserCacheDuration = getIntOrDefault(settingsMap, "browser_cache_duration", 7)
	settings.EnableCDN = getBoolOrDefault(settingsMap, "enable_cdn", false)
	settings.CDNURL = settingsMap["cdn_url"]
	settings.CDNImages = getBoolOrDefault(settingsMap, "cdn_images", true)
	settings.CDNCSS = getBoolOrDefault(settingsMap, "cdn_css", true)
	settings.CDNJS = getBoolOrDefault(settingsMap, "cdn_js", true)
	settings.CDNFonts = getBoolOrDefault(settingsMap, "cdn_fonts", true)
	settings.EnableImageOptimization = getBoolOrDefault(settingsMap, "enable_image_optimization", true)
	settings.ImageQuality = getIntOrDefault(settingsMap, "image_quality", 80)
	settings.EnableWebP = getBoolOrDefault(settingsMap, "enable_webp", true)
	settings.EnableLazyLoading = getBoolOrDefault(settingsMap, "enable_lazy_loading", true)
	settings.EnableJSMinify = getBoolOrDefault(settingsMap, "enable_js_minify", true)
	settings.EnableCSSMinify = getBoolOrDefault(settingsMap, "enable_css_minify", true)
	settings.EnableHTMLMinify = getBoolOrDefault(settingsMap, "enable_html_minify", false)
	settings.CombineJSFiles = getBoolOrDefault(settingsMap, "combine_js_files", true)
	settings.CombineCSSFiles = getBoolOrDefault(settingsMap, "combine_css_files", true)

	// Yedekleme Ayarları
	settings.EnableAutoBackup = getBoolOrDefault(settingsMap, "enable_auto_backup", false)
	settings.BackupFrequency = getOrDefault(settingsMap, "backup_frequency", "daily")
	settings.BackupDay = getOrDefault(settingsMap, "backup_day", "1")
	settings.BackupDate = getOrDefault(settingsMap, "backup_date", "1")
	settings.BackupTime = getOrDefault(settingsMap, "backup_time", "02:00")
	settings.BackupDatabase = getBoolOrDefault(settingsMap, "backup_database", true)
	settings.BackupFiles = getBoolOrDefault(settingsMap, "backup_files", true)
	settings.BackupSettings = getBoolOrDefault(settingsMap, "backup_settings", true)
	settings.BackupStorage = getOrDefault(settingsMap, "backup_storage", "local")
	settings.LocalBackupPath = getOrDefault(settingsMap, "local_backup_path", "./backups")
	settings.CloudAPIKey = settingsMap["cloud_api_key"]
	settings.CloudSecret = settingsMap["cloud_secret"]
	settings.CloudFolder = getOrDefault(settingsMap, "cloud_folder", "backups")
	settings.BackupRetention = getIntOrDefault(settingsMap, "backup_retention", 5)
	settings.CompressBackups = getBoolOrDefault(settingsMap, "compress_backups", true)
	settings.EncryptBackups = getBoolOrDefault(settingsMap, "encrypt_backups", false)
	settings.BackupNotifications = getBoolOrDefault(settingsMap, "backup_notifications", false)
	settings.BackupNotificationEmail = settingsMap["backup_notification_email"]
}

func getIntOrDefault(m map[string]string, key string, defaultValue int) int {
	if val, exists := m[key]; exists && val != "" {
		if intVal, err := json.Number(val).Int64(); err == nil {
			return int(intVal)
		}
	}
	return defaultValue
}
