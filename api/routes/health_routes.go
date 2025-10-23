package routes

import (
	"api-shiners/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func HealthRoutes(app *fiber.App, healthController handlers.HealthController) {
	api := app.Group("/api")
	api.Get("/health/database", healthController.HealthCheckDatabase)

	api.Get("/health/redis", healthController.HealthCheckRedis)
}


