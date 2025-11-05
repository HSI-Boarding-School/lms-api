package handlers

import (
	"context"
	"net/http"

	"api-shiners/pkg/user"
	"api-shiners/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController struct {
	userService user.UserService
}

func NewUserController(userService user.UserService) *UserController {
	return &UserController{userService}
}

func (ctrl *UserController) GetAllUsers(c *fiber.Ctx) error {
	users, err := ctrl.userService.GetAllUsers()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error(), "InternalServerError", nil)
	}

	return utils.Success(c, http.StatusOK, "Get all users successfully", users, nil)
}


func (ctrl *UserController) GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid user ID format", "InvalidUUID", nil)
	}

	user, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "User not found", "UserNotFound", nil)
	}

	return utils.Success(c, http.StatusOK, "Get user by ID successfully", user, nil)
}

func (ctrl *UserController) SetUserRole(c *fiber.Ctx) error {
	ctx := context.Background()

	// 🔹 Ambil user_id dari URL
	userIDParam := c.Params("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid user ID format", "InvalidUUID", nil)
	}

	// 🔹 Ambil nama role dari body
	var req struct {
		Role string `json:"role"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request body", "BadRequestException", nil)
	}

	if req.Role == "" {
		return utils.Error(c, http.StatusBadRequest, "Role name is required", "BadRequestException", nil)
	}

	// 🔹 Jalankan service → ambil user hasil update
	updatedUser, err := ctrl.userService.SetUserRole(ctx, userID, req.Role)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error(), "SetRoleException", nil)
	}

	// 🔹 Siapkan response data
	data := fiber.Map{
		"user_id": updatedUser.ID,
		"name":    updatedUser.Name,
		"email":   updatedUser.Email,
		"role":    req.Role,
	}

	// 🔹 Response sukses
	return utils.Success(c, http.StatusOK, "User role assigned successfully", data, nil)
}


