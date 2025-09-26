package middlewares

import (
	"strings"

	"github.com/daffa-fawwaz/shiners-lms-backend/app/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware untuk proteksi route
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ambil token dari header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// cek format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		tokenString := parts[1]

		// validasi token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// simpan userId ke context, biar bisa dipakai di handler berikutnya
		c.Locals("userId", claims.UserID)

		return c.Next()
	}
}
