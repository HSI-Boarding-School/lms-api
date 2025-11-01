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

	api.Post("/users/:id/role", middleware.AdminMiddleware, userController.SetUserRole)

	api.Post("/users/:id/deactivate", middleware.AdminMiddleware, userController.DeactivateUser)
	api.Post("/users/:id/activate", middleware.AdminMiddleware, userController.ActivateUser)

	api.Get("/profile", middleware.AuthMiddleware, userController.Profile)
	api.Put("/profile", middleware.AuthMiddleware, userController.UpdateProfile)
}