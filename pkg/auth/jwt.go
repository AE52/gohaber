package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/username/haber/internal/config"
	"github.com/username/haber/internal/domain"
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
func NewJWTAuth(secret string, accessTokenExp, refreshTokenExp int) *JWTAuth {
	cfg := &config.JWTConfig{
		Secret:          secret,
		AccessTokenExp:  accessTokenExp,
		RefreshTokenExp: refreshTokenExp,
	}
	return &JWTAuth{cfg: cfg}
}

// GenerateTokens erişim ve yenileme tokenlarını oluşturur
func (j *JWTAuth) GenerateTokens(user *domain.User) (accessToken, refreshToken string, err error) {
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
func (j *JWTAuth) generateToken(user *domain.User, tokenType TokenType, expiration time.Duration) (string, error) {
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

// ValidateToken tokeni doğrular ve kullanıcı bilgilerini döndürür
func (j *JWTAuth) ValidateToken(tokenString string) (*domain.User, error) {
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

	// Token tipini kontrol et (sadece erişim tokenleri için)
	if claims.TokenType != string(AccessToken) {
		return nil, errors.New("token tipini doğrulama hatası, access token bekleniyor")
	}

	// Kullanıcı nesnesini oluştur ve döndür
	user := &domain.User{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
	}

	return user, nil
}

// ValidateRefreshToken yenileme tokenini doğrular ve kullanıcı bilgilerini döndürür
func (j *JWTAuth) ValidateRefreshToken(refreshTokenString string) (*domain.User, error) {
	// Token'ı parse et
	token, err := jwt.ParseWithClaims(refreshTokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	// Token tipini kontrol et
	if claims.TokenType != string(RefreshToken) {
		return nil, errors.New("token tipini doğrulama hatası, refresh token bekleniyor")
	}

	// Kullanıcı nesnesini oluştur ve döndür
	user := &domain.User{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
	}

	return user, nil
}
