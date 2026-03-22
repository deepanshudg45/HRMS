// internal/auth/service.go
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secret = []byte("secret")

func GenerateAccessToken(empID string) (string, error) {
	claims := jwt.MapClaims{
		"employee_id": empID,
		"exp":         time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GenerateRefreshToken(empID string) (string, error) {
	claims := jwt.MapClaims{
		"employee_id": empID,
		"exp":         time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func HashToken(token string) (string, error) {
	return string(bcrypt.GenerateFromPassword([]byte(token), 12))
}

func VerifyTokenHash(hash, token string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(token)) == nil
}