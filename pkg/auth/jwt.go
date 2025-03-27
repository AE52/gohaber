package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/username/haber/internal/config"
	"github.com/username/haber/internal/models"
)

// TokenType token tipini belirtmek için kullanılır
type TokenType string

const (
	// AccessToken erişim token'ı
	AccessToken TokenType = "access"
	// RefreshToken yenileme token'ı
	RefreshToken TokenType = "refresh"
)

// JWTCustomClaims jwt için özel alanlar
type JWTCustomClaims struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// JWTAuth JWT kimlik doğrulama işlemlerini yönetir
type JWTAuth struct {
	cfg *config.JWTConfig
}

// NewJWTAuth yeni bir JWTAuth örneği oluşturur
func NewJWTAuth(cfg *config.JWTConfig) *JWTAuth {
	return &JWTAuth{cfg: cfg}
}

// GenerateTokens erişim ve yenileme tokenlarını oluşturur
func (j *JWTAuth) GenerateTokens(user *models.User) (accessToken, refreshToken string, err error) {
	// Access token oluşturma
	accessToken, err = j.generateToken(user, AccessToken, time.Minute*time.Duration(j.cfg.AccessTokenExp))
	if err != nil {
		return "", "", fmt.Errorf("erişim tokeni oluşturulurken hata: %w", err)
	}

	// Refresh token oluşturma
	refreshToken, err = j.generateToken(user, RefreshToken, time.Hour*time.Duration(j.cfg.RefreshTokenExp))
	if err != nil {
		return "", "", fmt.Errorf("yenileme tokeni oluşturulurken hata: %w", err)
	}

	return accessToken, refreshToken, nil
}

// generateToken belirtilen tipte ve sürede token oluşturur
func (j *JWTAuth) generateToken(user *models.User, tokenType TokenType, expiration time.Duration) (string, error) {
	// Token sona erme süresi
	expirationTime := time.Now().Add(expiration)

	// Token claim'leri
	claims := &JWTCustomClaims{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: string(tokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			// Standart JWT claim'leri
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "haber-sitesi",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	// Token oluşturma
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Token'ı imzalayıp dizeye dönüştür
	tokenString, err := token.SignedString([]byte(j.cfg.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken tokeni doğrular ve claim'leri döndürür
func (j *JWTAuth) ValidateToken(tokenString string) (*JWTCustomClaims, error) {
	// Token'ı parse et
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// HS256 algoritmasını kullandığımızı kontrol et
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("beklenmeyen imza metodu: %v", token.Header["alg"])
		}

		// Secret key döndür
		return []byte(j.cfg.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parse hatası: %w", err)
	}

	// Token'ın geçerli olup olmadığını kontrol et
	if !token.Valid {
		return nil, errors.New("token geçerli değil")
	}

	// Token claim'lerini döndür
	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok {
		return nil, errors.New("token claim'leri alınamadı")
	}

	return claims, nil
}

// RefreshAccessToken yenileme tokeni ile yeni bir erişim tokeni oluşturur
func (j *JWTAuth) RefreshAccessToken(refreshTokenString string) (string, error) {
	// Refresh token'ı doğrula
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		return "", fmt.Errorf("yenileme tokeni doğrulanamadı: %w", err)
	}

	// Token tipini kontrol et
	if claims.TokenType != string(RefreshToken) {
		return "", errors.New("token tipini doğrulama hatası, refresh token bekleniyor")
	}

	// Kullanıcı nesnesini oluştur
	user := &models.User{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
	}

	// Yeni erişim tokeni oluştur
	accessToken, err := j.generateToken(user, AccessToken, time.Minute*time.Duration(j.cfg.AccessTokenExp))
	if err != nil {
		return "", fmt.Errorf("yeni erişim tokeni oluşturulurken hata: %w", err)
	}

	return accessToken, nil
}
