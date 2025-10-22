package handlers

import (
	"api-shiners/pkg/auth"
	"api-shiners/pkg/utils"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService auth.AuthService
}

func NewAuthController(authService auth.AuthService) AuthController {
	return AuthController{authService: authService}
}

// ✅ Register
func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req auth.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request body", "BadRequestException", []utils.FieldError{
			{Field: "body", Messages: []string{"Invalid JSON format"}, Message: "Invalid JSON format"},
		})
	}

	createdUser, err := ctrl.authService.Register(context.Background(), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error(), "BadRequestException", []utils.FieldError{
			{Field: "email", Messages: []string{err.Error()}, Message: err.Error()},
		})
	}

	return utils.Success(c, http.StatusCreated, "User registered successfully", createdUser, nil)
}

// ✅ Login
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req auth.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request body", "BadRequestException", nil)
	}

	token, err := ctrl.authService.Login(context.Background(), req)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err.Error(), "UnauthorizedException", nil)
	}

	data := fiber.Map{"token": token}
	return utils.Success(c, http.StatusOK, "Login successful", data, nil)
}

func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return utils.Error(c, http.StatusBadRequest, "Authorization token required", "BadRequestException", nil)
	}

	// Hapus prefix "Bearer "
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	err := ctrl.authService.Logout(context.Background(), token)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error(), "LogoutException", nil)
	}

	return utils.Success(c, http.StatusOK, "Logout successful", nil, nil)
}

func (ctrl *AuthController) ForgotPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil || req.Email == "" {
		return utils.Error(c, http.StatusBadRequest, "Email is required", "BadRequestException", []utils.FieldError{
			{Field: "email", Messages: []string{"Email is required"}, Message: "Email is required"},
		})
	}

	token, err := ctrl.authService.GenerateResetToken(context.Background(), req.Email)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error(), "BadRequestException", nil)
	}

	return utils.Success(c, http.StatusOK, "Password reset token generated", fiber.Map{
		"email": req.Email,
		"token": token, // tampilkan untuk testing (di production sebaiknya kirim via email)
	}, nil)
}

func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil || req.Token == "" || req.NewPassword == "" {
		return utils.Error(c, http.StatusBadRequest, "Token and new password required", "BadRequestException", []utils.FieldError{
			{Field: "token", Messages: []string{"Token is required"}, Message: "Token is required"},
			{Field: "new_password", Messages: []string{"New password is required"}, Message: "New password is required"},
		})
	}

	if err := ctrl.authService.ResetPassword(context.Background(), req.Token, req.NewPassword); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error(), "BadRequestException", nil)
	}

	return utils.Success(c, http.StatusOK, "Password has been reset successfully", nil, nil)
}


