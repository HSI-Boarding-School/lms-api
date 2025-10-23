package middleware

import (
	"net/http"
	"os"
	"strings"

	"api-shiners/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware memastikan user sudah login dan token valid
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.Error(c, http.StatusUnauthorized, "Missing authorization header", "UnauthorizedException", nil)
	}

	// Ambil token dari header: "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return utils.Error(c, http.StatusUnauthorized, "Invalid token format", "UnauthorizedException", nil)
	}

	// Parse token JWT
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(http.StatusUnauthorized, "Invalid token signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return utils.Error(c, http.StatusUnauthorized, "Invalid or expired token", "UnauthorizedException", nil)
	}

	// Ambil claims dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utils.Error(c, http.StatusUnauthorized, "Invalid token claims", "UnauthorizedException", nil)
	}

	// Ambil user_id dari claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return utils.Error(c, http.StatusUnauthorized, "Invalid user ID in token", "UnauthorizedException", nil)
	}

	// Validasi format UUID
	_, err = uuid.Parse(userID)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, "Invalid user ID format", "UnauthorizedException", nil)
	}

	// Simpan user_id ke context untuk digunakan di handler
	c.Locals("user_id", userID)

	return c.Next()
}
