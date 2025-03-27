package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword parolayı şifreler
func HashPassword(password string) (string, error) {
	// Parola şifreleme işlemi
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword şifrelenmiş parola ile girilen parolayı karşılaştırır
func CheckPassword(hashedPassword, password string) bool {
	// Parolaları karşılaştır
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
