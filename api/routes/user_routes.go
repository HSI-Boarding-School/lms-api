package routes

import (
	"api-shiners/api/handlers"
	"api-shiners/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, userController *handlers.UserController) {
	api := app.Group("/api")
	api.Get("/users", middleware.AdminMiddleware, userController.GetAllUsers)
	api.Get("/users/:id", middleware.AuthMiddleware, userController.GetUserByID)
}