package handlers

import (
	"api-shiners/pkg/auth"
	"api-shiners/pkg/utils"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AuthController handles authentication related endpoints
type AuthController struct {
	authService auth.AuthService
}

func NewAuthController(authService auth.AuthService) AuthController {
	return AuthController{authService: authService}
}

// ==================== REGISTER ====================

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "Register Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
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

// ==================== LOGIN ====================

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login Request"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/auth/login [post]
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req auth.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request body", "BadRequestException", nil)
	}

	token, exp, err := ctrl.authService.Login(context.Background(), req)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err.Error(), "UnauthorizedException", nil)
	}

	data := fiber.Map{
		"token":       token,
		"expires_in":  exp.Format(time.RFC3339),
		"token_type":  "Bearer",
	}

	return utils.Success(c, http.StatusOK, "Login successful", data, nil)
}

// ==================== LOGOUT ====================

// Logout godoc
// @Summary Logout user
// @Description Invalidate user token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
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

// ==================== FORGOT PASSWORD ====================

// ForgotPassword godoc
// @Summary Request password reset
// @Description Generate reset token and send it to user's email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Email Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/auth/forgot-password [post]
func (ctrl *AuthController) ForgotPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil || req.Email == "" {
		return utils.Error(c, http.StatusBadRequest, "Email is required", "BadRequestException", nil)
	}

	token, err := ctrl.authService.GenerateResetToken(context.Background(), req.Email)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error(), "BadRequestException", nil)
	}

	return utils.Success(c, http.StatusOK, "Password reset token generated", fiber.Map{
		"email": req.Email,
		"token": token, // tampilkan untuk testing
	}, nil)
}

// ==================== RESET PASSWORD ====================

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset password using valid reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Reset Password Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/auth/reset-password [post]
func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil || req.Token == "" || req.NewPassword == "" {
		return utils.Error(c, http.StatusBadRequest, "Token and new password required", "BadRequestException", nil)
	}

	if err := ctrl.authService.ResetPassword(context.Background(), req.Token, req.NewPassword); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error(), "BadRequestException", nil)
	}

	return utils.Success(c, http.StatusOK, "Password has been reset successfully", nil, nil)
}
