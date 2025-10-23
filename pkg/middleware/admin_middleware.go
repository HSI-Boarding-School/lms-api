package middleware

import (
	"net/http"
	"os"
	"strings"

	"api-shiners/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AdminOnly middleware hanya mengizinkan akses untuk role ADMIN
func AdminMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.Error(c, http.StatusUnauthorized, "Missing authorization header", "UnauthorizedException", nil)
	}

	// Ambil token dari header: "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return utils.Error(c, http.StatusUnauthorized, "Invalid token format", "UnauthorizedException", nil)
	}

	// Parse token
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

	// Ambil claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utils.Error(c, http.StatusUnauthorized, "Invalid token claims", "UnauthorizedException", nil)
	}

	// Ambil roles dari claims
	roles, ok := claims["roles"].([]interface{})
	if !ok {
		return utils.Error(c, http.StatusForbidden, "You are not authorized", "ForbiddenException", nil)
	}

	// Cek apakah user memiliki role ADMIN
	isAdmin := false
	for _, r := range roles {
		if roleMap, ok := r.(map[string]interface{}); ok {
			if roleName, ok := roleMap["name"].(string); ok && roleName == "ADMIN" {
				isAdmin = true
				break
			}
		}
	}

	if !isAdmin {
		return utils.Error(c, http.StatusForbidden, "Access restricted to ADMIN only", "ForbiddenException", nil)
	}

	return c.Next()
}
