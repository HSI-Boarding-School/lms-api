// @title Shiners API Documentation
// @version 1.0
// @description This is the authentication API documentation for Shiners project.
// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"fmt"
	"log"
	"os"

	"api-shiners/api/handlers"
	"api-shiners/api/routes"
	"api-shiners/pkg/auth"
	"api-shiners/pkg/user"
	"api-shiners/pkg/config"

	_ "api-shiners/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	godotenv.Load()

	// Koneksi ke database
	config.ConnectDatabase()

	// Koneksi ke Redis
	config.InitRedis()

	// Buat Fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
	}))

	// Ambil port dari .env
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	authRepo := auth.NewUserRepository(config.DB)
	authService := auth.NewAuthService(authRepo)
	authController := handlers.NewAuthController(authService)

	healthController := handlers.NewHealthController()

	userRepo := user.NewUserRepository(config.DB)
	userService := user.NewUserService(userRepo)
	userController := handlers.NewUserController(userService)

	routes.UserRoutes(app, userController)
	routes.HealthRoutes(app, healthController)
	routes.AuthRoutes(app, authController)

	log.Printf("🚀 Server running on port %s...\n", port)
	app.Listen(fmt.Sprintf(":%s", port))
}
