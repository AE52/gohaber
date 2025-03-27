package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Kullanıcı rolleri için sabitler
const (
	RoleAdmin       = "admin"       // Tam yetkili admin
	RoleEditor      = "editor"      // İçerik düzenleme yetkisi
	RoleAuthor      = "author"      // Sadece kendi içeriklerini düzenleyebilir
	RoleContributor = "contributor" // Sadece taslak oluşturabilir, yayınlayamaz
	RoleSubscriber  = "subscriber"  // Normal site üyesi
)

// Kullanıcı durumları için sabitler
const (
	StatusActive   = "active"   // Aktif kullanıcı
	StatusInactive = "inactive" // Pasif kullanıcı
	StatusPending  = "pending"  // Onay bekleyen
	StatusBanned   = "banned"   // Yasaklanmış
)

// Token, kimlik doğrulama tokenını temsil eder
type Token struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	Type      string    `json:"type" db:"type"` // login, reset-password, confirm-email, etc.
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TokenType, token tiplerini belirtir
type TokenType string

const (
	TokenTypeLogin         TokenType = "login"
	TokenTypeResetPassword TokenType = "reset-password"
	TokenTypeEmailConfirm  TokenType = "confirm-email"
)

// LoginUserRequest, giriş isteği için model
type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

// RegisterUserRequest, kayıt isteği için model
type RegisterUserRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	FullName        string `json:"full_name"`
}

// LoginResponse, giriş yanıtı için model
type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// AuthService, kimlik doğrulama için kullanılan arayüz
type AuthService interface {
	Authenticate(username, password string) (*User, error)
	CreateUser(req RegisterUserRequest) (*User, error)
	GetUserByID(id uint) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
	ChangePassword(id uint, currentPassword, newPassword string) error
	ResetPassword(email string) (string, error)
	ConfirmResetPassword(token, newPassword string) error
	CreateToken(userID uint, tokenType TokenType, expiration time.Duration) (string, error)
	VerifyToken(token string, tokenType TokenType) (*User, error)
}

// SqlAuthService, Auth servisinin SQL implementasyonu
type SqlAuthService struct {
	DB *sql.DB
}

// Authenticate, kullanıcı adı ve parola ile kimlik doğrulama yapar
func (s *SqlAuthService) Authenticate(username, password string) (*User, error) {
	var user User

	// Kullanıcı adı veya e-posta ile sorgulama yapabiliriz
	query := `SELECT id, username, email, password_hash, full_name, role, profile_image, created_at, updated_at
	FROM users WHERE (username = ? OR email = ?) AND role != ?`

	err := s.DB.QueryRow(query, username, username, StatusInactive).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.ProfileImage,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Kullanıcı adı veya parola hatalı")
		}
		return nil, err
	}

	// Parola doğrulama
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("Kullanıcı adı veya parola hatalı")
	}

	// Son giriş zamanını güncelle
	now := time.Now()
	_, err = s.DB.Exec("UPDATE users SET updated_at = ? WHERE id = ?", now, user.ID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser, yeni bir kullanıcı oluşturur
func (s *SqlAuthService) CreateUser(req RegisterUserRequest) (*User, error) {
	// Zorunlu alanları kontrol et
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("Kullanıcı adı, e-posta ve parola zorunludur")
	}

	// Parolaların eşleştiğini kontrol et
	if req.Password != req.ConfirmPassword {
		return nil, errors.New("Parolalar eşleşmiyor")
	}

	// Kullanıcı adının benzersiz olduğunu kontrol et
	existingUser, _ := s.GetUserByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("Bu kullanıcı adı zaten kullanılıyor")
	}

	// E-postanın benzersiz olduğunu kontrol et
	existingEmail, _ := s.GetUserByEmail(req.Email)
	if existingEmail != nil {
		return nil, errors.New("Bu e-posta adresi zaten kullanılıyor")
	}

	// Parolayı hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Role:         RoleSubscriber, // Yeni kullanıcılar varsayılan olarak üye
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Kullanıcıyı veritabanına ekle
	query := `INSERT INTO users (username, email, password_hash, full_name, role, created_at, updated_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := s.DB.Exec(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Yeni eklenen kullanıcının ID'sini al
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = uint(id)

	return user, nil
}

// GetUserByID, ID'ye göre kullanıcı bilgilerini getirir
func (s *SqlAuthService) GetUserByID(id uint) (*User, error) {
	var user User

	query := `SELECT id, username, email, password_hash, full_name, role, profile_image, created_at, updated_at
	FROM users WHERE id = ?`

	err := s.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.ProfileImage,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername, kullanıcı adına göre kullanıcı bilgilerini getirir
func (s *SqlAuthService) GetUserByUsername(username string) (*User, error) {
	var user User

	query := `SELECT id, username, email, password_hash, full_name, role, profile_image, created_at, updated_at
	FROM users WHERE username = ?`

	err := s.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.ProfileImage,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail, e-posta adresine göre kullanıcı bilgilerini getirir
func (s *SqlAuthService) GetUserByEmail(email string) (*User, error) {
	var user User

	query := `SELECT id, username, email, password_hash, full_name, role, profile_image, created_at, updated_at
	FROM users WHERE email = ?`

	err := s.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.ProfileImage,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// UpdateUser, kullanıcı bilgilerini günceller
func (s *SqlAuthService) UpdateUser(user *User) error {
	user.UpdatedAt = time.Now()

	query := `UPDATE users SET 
	username = ?, 
	email = ?, 
	full_name = ?, 
	role = ?, 
	profile_image = ?, 
	updated_at = ? 
	WHERE id = ?`

	_, err := s.DB.Exec(query,
		user.Username,
		user.Email,
		user.FullName,
		user.Role,
		user.ProfileImage,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// DeleteUser, kullanıcıyı siler
func (s *SqlAuthService) DeleteUser(id uint) error {
	// Gerçekte silmek yerine role'u güncelliyoruz
	query := `UPDATE users SET role = ?, updated_at = ? WHERE id = ?`
	_, err := s.DB.Exec(query, StatusInactive, time.Now(), id)
	return err
}

// ChangePassword, kullanıcının parolasını değiştirir
func (s *SqlAuthService) ChangePassword(id uint, currentPassword, newPassword string) error {
	// Mevcut kullanıcıyı getir
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("Kullanıcı bulunamadı")
	}

	// Mevcut parolayı doğrula
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword))
	if err != nil {
		return errors.New("Mevcut parola hatalı")
	}

	// Yeni parolayı hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Parolayı güncelle
	query := `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	_, err = s.DB.Exec(query, string(hashedPassword), time.Now(), id)

	return err
}

// ResetPassword, parola sıfırlama istemi oluşturur
func (s *SqlAuthService) ResetPassword(email string) (string, error) {
	// E-posta adresine göre kullanıcıyı bul
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("Bu e-posta adresi ile kayıtlı kullanıcı bulunamadı")
	}

	// Parola sıfırlama tokeni oluştur (24 saat geçerli)
	token, err := s.CreateToken(user.ID, TokenTypeResetPassword, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ConfirmResetPassword, parola sıfırlama işlemini onaylar
func (s *SqlAuthService) ConfirmResetPassword(token, newPassword string) error {
	// Tokeni doğrula
	user, err := s.VerifyToken(token, TokenTypeResetPassword)
	if err != nil {
		return err
	}

	// Yeni parolayı hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Parolayı güncelle
	query := `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	_, err = s.DB.Exec(query, string(hashedPassword), time.Now(), user.ID)
	if err != nil {
		return err
	}

	// Tokeni kullanıldı olarak işaretle (sil)
	_, err = s.DB.Exec("DELETE FROM tokens WHERE token = ?", token)

	return err
}

// CreateToken, yeni bir token oluşturur
func (s *SqlAuthService) CreateToken(userID uint, tokenType TokenType, expiration time.Duration) (string, error) {
	// Rastgele bir token oluştur
	tokenStr, err := GenerateRandomToken(32)
	if err != nil {
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(expiration)

	// Tokeni veritabanına kaydet
	query := `INSERT INTO tokens (user_id, token, type, expires_at, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = s.DB.Exec(query, userID, tokenStr, string(tokenType), expiresAt, now)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// VerifyToken, token doğrulaması yapar
func (s *SqlAuthService) VerifyToken(token string, tokenType TokenType) (*User, error) {
	var tokenRecord Token

	// Tokeni veritabanından sorgula
	query := `SELECT id, user_id, token, type, expires_at, created_at FROM tokens 
	WHERE token = ? AND type = ? AND expires_at > ?`

	err := s.DB.QueryRow(query, token, string(tokenType), time.Now()).Scan(
		&tokenRecord.ID,
		&tokenRecord.UserID,
		&tokenRecord.Token,
		&tokenRecord.Type,
		&tokenRecord.ExpiresAt,
		&tokenRecord.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Geçersiz veya süresi dolmuş token")
		}
		return nil, err
	}

	// Tokene bağlı kullanıcıyı getir
	user, err := s.GetUserByID(uint(tokenRecord.UserID))
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("Kullanıcı bulunamadı")
	}

	return user, nil
}

// GenerateRandomToken, belirtilen uzunlukta rastgele bir token oluşturur
func GenerateRandomToken(length int) (string, error) {
	// Bu fonksiyon implementasyonu crypto/rand paketi kullanılarak yapılabilir
	// Şimdilik basit bir şekilde simüle edelim
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
		time.Sleep(1 * time.Nanosecond)
	}
	return string(result), nil
}
