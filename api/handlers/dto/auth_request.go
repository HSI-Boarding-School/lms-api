package dto

type RegisterRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"strongpassword123"`
}

type RegisterResponse struct {
	ID    string `json:"id" example:"a3b2c1d4-56ef-7890-gh12-ijk345lmn678"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"strongpassword123"`
}

type LoginResponse struct {
	Token      string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	ExpiresIn  string      `json:"expires_in" example:"2025-10-18T15:04:05Z"`
	TokenType  string      `json:"token_type" example:"Bearer"`
	User       interface{} `json:"user"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" example:"john@example.com"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" example:"123456"`
	NewPassword string `json:"new_password" example:"newStrongPassword123"`
}

type GenericResponse struct {
	Message string `json:"message" example:"Operation successful"`
}
