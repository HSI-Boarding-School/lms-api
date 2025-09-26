// controllers/auth_controller.go
package controllers

import (
	"github.com/daffa-fawwaz/shiners-lms-backend/repositories"
	"github.com/daffa-fawwaz/shiners-lms-backend/services"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService services.AuthService
	userRepo    repositories.UserRepository
}

func NewAuthController(authService services.AuthService, userRepo repositories.UserRepository) 			AuthController {
	return AuthController{
		authService: authService,
		userRepo:    userRepo,
	}
}

func (ac AuthController) Register(c *fiber.Ctx) error {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := ac.authService.Register(c.Context(), body.Name, body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (ac AuthController) Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, err := ac.userRepo.FindByEmail(c.Context(), body.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	accessToken, refreshToken, role, err := ac.authService.Login(c.Context(), body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"info": fiber.Map{
			"role":  role,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}




