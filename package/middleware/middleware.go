package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type JWTClaims struct {
	EmployeeID string `json:"employeeId"`
	Role       string `json:"role"`
	ExpiresAt  int64  `json:"exp"`
}

func GenerateJWT(employeeID string, role string, secret string) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	claims := JWTClaims{
		EmployeeID: strings.TrimSpace(employeeID),
		Role:       strings.ToUpper(strings.TrimSpace(role)),
		ExpiresAt:  time.Now().Add(24 * time.Hour).Unix(),
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signingInput := headerEncoded + "." + payloadEncoded
	signature := signJWT(signingInput, secret)

	return signingInput + "." + signature, nil
}

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := strings.TrimSpace(c.Get("Authorization"))
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing bearer token",
			})
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := ParseJWT(token, os.Getenv("JWT_SECRET"))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		c.Locals("employeeId", claims.EmployeeID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func CurrentEmployeeID(c *fiber.Ctx) string {
	if employeeID, ok := c.Locals("employeeId").(string); ok {
		return strings.TrimSpace(employeeID)
	}

	return ""
}

func CurrentRole(c *fiber.Ctx) string {
	if role, ok := c.Locals("role").(string); ok {
		return strings.ToUpper(strings.TrimSpace(role))
	}

	return ""
}

func IsHR(c *fiber.Ctx) bool {
	return CurrentRole(c) == "HR"
}

func CanAccessEmployee(c *fiber.Ctx, employeeID string) bool {
	employeeID = strings.TrimSpace(employeeID)
	if employeeID == "" {
		return false
	}

	if IsHR(c) {
		return true
	}

	return CurrentEmployeeID(c) == employeeID
}

func ParseJWT(token string, secret string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fiber.ErrUnauthorized
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSignature := signJWT(signingInput, secret)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return nil, fiber.ErrUnauthorized
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, err
	}

	if claims.ExpiresAt != 0 && time.Now().Unix() > claims.ExpiresAt {
		return nil, fiber.ErrUnauthorized
	}

	return &claims, nil
}

func signJWT(payload string, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}
