package app

import (
	"log"

	"github.com/daffa-fawwaz/shiners-lms-backend/app/controllers"
	"github.com/daffa-fawwaz/shiners-lms-backend/app/repositories"
	"github.com/daffa-fawwaz/shiners-lms-backend/app/routes"
	"github.com/daffa-fawwaz/shiners-lms-backend/app/services"
	"github.com/daffa-fawwaz/shiners-lms-backend/config"

	"github.com/gofiber/fiber/v2"
)

func SetupApp() *fiber.App {
    db, err := config.ConnectDB()
    if err != nil {
        log.Fatal("failed to connect database: ", err)
    }

    userRepo := repositories.NewUserRepository(db)
    roleRepo := repositories.NewRoleRepository(db)

    authService := services.NewAuthService(userRepo, roleRepo)

	authController := controllers.NewAuthController(authService, userRepo)

    app := fiber.New()
    routes.AuthRoutes(app, authController)

    return app
}
