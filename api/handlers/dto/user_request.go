package dto


type SetRoleRequest struct {
	Role string `json:"role" example:"ADMIN"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
	Role     string `json:"role,omitempty"`
}

type UserRoleResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserStatusResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
}

type UserProfileResponse struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	IsActive bool     `json:"is_active"`
	Roles    []string `json:"roles"`
}

type MetaResponse struct {
	Page    int `json:"page" example:"1"`
	PerPage int `json:"per_page" example:"10"`
	Total   int `json:"total" example:"100"`
}


type PaginatedUsersResponse struct {
	Data []UserResponse `json:"data"`
	Meta MetaResponse   `json:"meta"`
}
