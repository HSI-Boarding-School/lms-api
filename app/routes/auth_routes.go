package routes

import (
	"github.com/daffa-fawwaz/shiners-lms-backend/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authController controllers.AuthController) {
	auth := app.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
}

