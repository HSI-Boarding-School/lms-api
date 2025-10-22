package main

import (
	"fmt"
	"log"
	"os"

	"api-shiners/api/handlers"
	"api-shiners/api/routes"
	"api-shiners/pkg/database"
	"api-shiners/pkg/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	godotenv.Load()

	// Koneksi ke database
	database.ConnectDatabase()

	// Buat Fiber app
	app := fiber.New()

	// Ambil port dari .env
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	authRepo := auth.NewUserRepository(database.DB)
	authService := auth.NewAuthService(authRepo)

	authController := handlers.NewAuthController(authService)
	routes.AuthRoutes(app, authController)

	log.Printf("🚀 Server running on port %s...\n", port)
	app.Listen(fmt.Sprintf(":%s", port))
}
