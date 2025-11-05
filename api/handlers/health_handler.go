package handlers

import (
	"api-shiners/pkg/config"
	"api-shiners/pkg/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type HealthController struct{}

func NewHealthController() HealthController {
	return HealthController{}
}

// HealthCheck godoc
// @Summary Check service health
// @Description Check database connection status
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /api/health/database [get]
func (ctrl *HealthController) HealthCheckDatabase(c *fiber.Ctx) error {
	status := fiber.Map{
		"status": "ok",
	}

	sqlDB, err := config.DB.DB()
	if err != nil {
		status["database"] = "error: cannot access sql.DB"
		status["status"] = "degraded"
		return utils.Error(c, http.StatusInternalServerError, err.Error(), "InternalServerError", nil)
	}

	if err := sqlDB.Ping(); err != nil {
		status["database"] = "unreachable"
		status["status"] = "degraded"
		return utils.Error(c, http.StatusInternalServerError, err.Error(), "InternalServerError", nil)
	}

	status["database"] = "connected"
	return utils.Success(c, http.StatusOK, "Database connected successfully", status, nil)
}



// HealthCheck godoc
// @Summary Check service health
// @Description Check database connection status
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /api/health/redis [get]
func (ctrl *HealthController) HealthCheckRedis(c *fiber.Ctx) error {
	status := fiber.Map{
		"status": "ok",
	}

	ctx := context.Background()

	if config.RedisClient == nil {
		status["redis"] = "not connected"
		return utils.Success(c, http.StatusOK, "Redis not initialized (dev mode or disabled)", status, nil)
	}

	_, err := config.RedisClient.Ping(ctx).Result()
	if err != nil {
		status["redis"] = "not connected"
		return utils.Success(c, http.StatusOK, fmt.Sprintf("Redis not connected: %v", err.Error()), status, nil)
	}

	status["redis"] = "connected"
	return utils.Success(c, http.StatusOK, "Redis connected successfully", status, nil)
}

