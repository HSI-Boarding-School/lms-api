package handlers

import (
	"api-shiners/api/handlers/dto"
	"api-shiners/pkg/user"
	"api-shiners/pkg/utils"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController struct {
	userService user.UserService
}

func NewUserController(userService user.UserService) *UserController {
	return &UserController{userService: userService}
}


// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve a paginated list of all users
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} dto.PaginatedUsersResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/users [get]
func (ctrl *UserController) GetAllUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 10)

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	users, total, err := ctrl.userService.GetAllUsers(page, perPage)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error(), "InternalServerError", nil)
	}

	// mapping ke DTO
	var userDTOs []dto.UserResponse
	for _, u := range users {
		userDTOs = append(userDTOs, dto.UserResponse{
			ID:       u.ID.String(),
			Name:     u.Name,
			Email:    u.Email,
			IsActive: u.IsActive,
		})
	}

	meta := dto.MetaResponse{
		Page:    page,
		PerPage: perPage,
		Total:   int(total),
	}

	response := dto.PaginatedUsersResponse{
		Data: userDTOs,
		Meta: meta,
	}

	return utils.Success(c, http.StatusOK, "Get all users successfully", response, nil)
}


// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve user details by user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/users/{id} [get]
func (ctrl *UserController) GetUserByID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid user ID format", "InvalidUUID", nil)
	}

	u, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "User not found", "UserNotFound", nil)
	}

	resp := dto.UserResponse{
		ID:       u.ID.String(),
		Name:     u.Name,
		Email:    u.Email,
		IsActive: u.IsActive,
	}

	return utils.Success(c, http.StatusOK, "Get user by ID successfully", resp, nil)
}


// SetUserRole godoc
// @Summary Set user role
// @Description Assign or update a user's role
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body dto.SetRoleRequest true "Set Role Request"
// @Success 200 {object} dto.UserRoleResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/users/{id}/role [put]
func (ctrl *UserController) SetUserRole(c *fiber.Ctx) error {
	var req dto.SetRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request body", "BadRequestException", nil)
	}

	if req.Role == "" {
		return utils.Error(c, http.StatusBadRequest, "Role is required", "BadRequestException", nil)
	}

	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid user ID format", "InvalidUUID", nil)
	}

	updatedUser, err := ctrl.userService.SetUserRole(context.Background(), userID, req.Role)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error(), "SetRoleException", nil)
	}

	resp := dto.UserRoleResponse{
		ID:    updatedUser.ID.String(),
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
		Role:  req.Role,
	}

	return utils.Success(c, http.StatusOK, "User role assigned successfully", resp, nil)
}


// DeactivateUser godoc
// @Summary Deactivate user
// @Description Deactivate a user's account by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserStatusResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/users/{id}/deactivate [put]
func (ctrl *UserController) DeactivateUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid user ID format", "InvalidUUID", nil)
	}

	deactivatedUser, err := ctrl.userService.DeactivateUser(context.Background(), userID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error(), "DeactivateUserException", nil)
	}

	resp := dto.UserStatusResponse{
		ID:       deactivatedUser.ID.String(),
		Name:     deactivatedUser.Name,
		Email:    deactivatedUser.Email,
		IsActive: deactivatedUser.IsActive,
	}

	return utils.Success(c, http.StatusOK, "User deactivated successfully", resp, nil)
}


// ActivateUser godoc
// @Summary Activate user
// @Description Activate a previously deactivated user account
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserStatusResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/users/{id}/activate [put]
func (ctrl *UserController) ActivateUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid user ID format", "InvalidUUID", nil)
	}

	activatedUser, err := ctrl.userService.ActivateUser(context.Background(), userID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error(), "ActivateUserException", nil)
	}

	resp := dto.UserStatusResponse{
		ID:       activatedUser.ID.String(),
		Name:     activatedUser.Name,
		Email:    activatedUser.Email,
		IsActive: activatedUser.IsActive,
	}

	return utils.Success(c, http.StatusOK, "User activated successfully", resp, nil)
}


// Profile godoc
// @Summary Get user profile
// @Description Retrieve authenticated user's profile
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserProfileResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/users/profile [get]
func (ctrl *UserController) Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.Error(c, http.StatusUnauthorized, "Unauthorized", "UnauthorizedException", nil)
	}

	profile, err := ctrl.userService.GetProfile(context.Background(), userID.(string))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "User not found", "NotFoundException", nil)
	}

	var roles []string
	for _, r := range profile.Roles {
		roles = append(roles, string(r.Name))
	}

	resp := dto.UserProfileResponse{
		ID:       profile.ID.String(),
		Name:     profile.Name,
		Email:    profile.Email,
		IsActive: profile.IsActive,
		Roles:    roles,
	}

	return utils.Success(c, http.StatusOK, "Profile fetched successfully", resp, nil)
}

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserController) UpdateProfile(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return utils.Error(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized", nil)
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, "Invalid user ID", "Unauthorized", nil)
	}

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request body", "BadRequest", nil)
	}

	ctx := context.Background()
	updatedUser, err := h.userService.UpdateProfile(ctx, userID, req.Name, req.Email)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error(), "InternalServerError", nil)
	}

	return utils.Success(c, http.StatusOK, "Profile updated successfully", updatedUser, nil)
}
