package routes

import (
	"api-shiners/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authController handlers.AuthController) {
	api := app.Group("/api/auth")

	api.Post("/register", authController.Register)
	api.Post("/login", authController.Login)
	api.Post("/logout", authController.Logout)
	api.Post("/forgot-password", authController.ForgotPassword)
	api.Post("/reset-password", authController.ResetPassword)
}
