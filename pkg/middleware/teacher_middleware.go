package middleware

import (
	"api-shiners/pkg/utils"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// TeacherMiddleware memastikan hanya TEACHER yang bisa mengakses route
func TeacherMiddleware(c *fiber.Ctx) error {
	// Ambil header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.Error(c, http.StatusUnauthorized, "Missing Authorization header", "UnauthorizedException", nil)
	}

	// Ambil token dari header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return utils.Error(c, http.StatusUnauthorized, "Invalid token format", "UnauthorizedException", nil)
	}

	// Parse JWT token
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

	// Ambil role dari token (bisa string atau array)
	rawRoles, exists := claims["roles"]
	if !exists {
		return utils.Error(c, http.StatusForbidden, "No roles found in token", "ForbiddenException", nil)
	}

	isTeacher := false

	switch roles := rawRoles.(type) {
	case string:
		if strings.EqualFold(roles, "TEACHER") {
			isTeacher = true
		}
	case []interface{}:
		for _, r := range roles {
			if rs, ok := r.(string); ok && strings.EqualFold(rs, "TEACHER") {
				isTeacher = true
				break
			}
		}
	default:
		roleStr := strings.ToUpper(strings.TrimSpace(fmt.Sprint(rawRoles)))
		if strings.Contains(roleStr, "TEACHER") {
			isTeacher = true
		}
	}

	if !isTeacher {
		return utils.Error(c, http.StatusForbidden, "Access restricted to TEACHER only", "ForbiddenException", nil)
	}

	// Simpan info user ke context
	c.Locals("user_id", claims["user_id"])
	c.Locals("email", claims["email"])
	c.Locals("roles", rawRoles)

	return c.Next()
}
