package handlers

import (
	"api-shiners/api/handlers/dto"
	"api-shiners/pkg/auth"
	"api-shiners/pkg/utils"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService auth.AuthService
}

func NewAuthController(authService auth.AuthService) AuthController {
	return AuthController{authService: authService}
}


// @Summary Register a new user
// @Description Membuat akun user baru
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/auth/register [post]
func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req auth.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request body", "BadRequestException", nil)
	}

	createdUser, err := ctrl.authService.Register(context.Background(), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error(), "BadRequestException", nil)
	}

	return utils.Success(c, http.StatusCreated, "User registered successfully", createdUser, nil)
}


// @Summary Login user
// @Description Autentikasi user dan mendapatkan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /api/auth/login [post]
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req auth.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request body", "BadRequestException", nil)
	}

	user, token, exp, permissions, err := ctrl.authService.LoginCore(context.Background(), req)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err.Error(), "UnauthorizedException", nil)
	}

	data := fiber.Map{
		"token":      token,
		"expires_in": exp.Format(time.RFC3339),
		"token_type": "Bearer",
		"user": fiber.Map{
			"id":          user.ID,
			"name":        user.Name,
			"role":        user.Roles,
			"permissions": permissions,
		},
	}

	return utils.Success(c, http.StatusOK, "Login successful", data, nil)
}


// @Summary Logout user
// @Description Mengakhiri sesi dan menonaktifkan token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.GenericResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /api/auth/logout [post]
func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return utils.Error(c, http.StatusBadRequest, "Authorization token required", "BadRequestException", nil)
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	err := ctrl.authService.Logout(context.Background(), token)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error(), "LogoutException", nil)
	}

	return utils.Success(c, http.StatusOK, "Logout successful", nil, nil)
}


// @Summary Request password reset
// @Description Generate reset token dan kirim ke email user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Forgot Password Request"
// @Success 200 {object} dto.GenericResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/auth/forgot-password [post]
func (ctrl *AuthController) ForgotPassword(c *fiber.Ctx) error {
	var req dto.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil || req.Email == "" {
		return utils.Error(c, http.StatusBadRequest, "Email is required", "BadRequestException", nil)
	}

	token, err := ctrl.authService.GenerateResetToken(context.Background(), req.Email)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error(), "BadRequestException", nil)
	}

	return utils.Success(c, http.StatusOK, "Password reset token generated", fiber.Map{
		"email": req.Email,
		"token": token, // ditampilkan untuk keperluan testing
	}, nil)
}


// @Summary Reset user password
// @Description Reset password menggunakan reset token yang valid
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} dto.GenericResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/auth/reset-password [post]
func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil || req.Token == "" || req.NewPassword == "" {
		return utils.Error(c, http.StatusBadRequest, "Token and new password required", "BadRequestException", nil)
	}

	if err := ctrl.authService.ResetPassword(context.Background(), req.Token, req.NewPassword); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error(), "BadRequestException", nil)
	}

	return utils.Success(c, http.StatusOK, "Password has been reset successfully", nil, nil)
}
