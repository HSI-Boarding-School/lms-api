package handlers

import (
	"api-shiners/pkg/config"
	"api-shiners/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// HealthController untuk handle health check API
type HealthController struct{}

// NewHealthController membuat instance baru HealthController
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
// @Router /api/health [get]
func (ctrl *HealthController) HealthCheckDatabase(c *fiber.Ctx) error {
	status := fiber.Map{
		"status": "ok",
	}

	// ✅ Cek koneksi database PostgreSQL
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


func (ctrl *HealthController) HealthCheckRedis(c *fiber.Ctx) error {
	// Coba ping ke Redis
	_, err := config.RedisClient.Ping(config.Ctx).Result()
	if err != nil {
		// Jika gagal konek ke Redis
		return utils.Error(c, http.StatusInternalServerError, err.Error(), "InternalServerError", nil)
	}

	status := fiber.Map{
		"status": "ok",
	}

	status["redis"] = "connected"

	// Jika Redis sehat
	return utils.Success(c, http.StatusOK, "Redis connected successfully", status, nil)
}

