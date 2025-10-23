package handlers

import (
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
