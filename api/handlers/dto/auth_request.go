package dto

// RegisterRequest — payload untuk registrasi user
type RegisterRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"strongpassword123"`
}

// RegisterResponse — response sukses registrasi
type RegisterResponse struct {
	ID    string `json:"id" example:"a3b2c1d4-56ef-7890-gh12-ijk345lmn678"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

// LoginRequest — payload login user
type LoginRequest struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"strongpassword123"`
}

// LoginResponse — response sukses login
type LoginResponse struct {
	Token      string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	ExpiresIn  string      `json:"expires_in" example:"2025-10-18T15:04:05Z"`
	TokenType  string      `json:"token_type" example:"Bearer"`
	User       interface{} `json:"user"`
}

// ForgotPasswordRequest — payload untuk lupa password
type ForgotPasswordRequest struct {
	Email string `json:"email" example:"john@example.com"`
}

// ResetPasswordRequest — payload untuk reset password
type ResetPasswordRequest struct {
	Token       string `json:"token" example:"123456"`
	NewPassword string `json:"new_password" example:"newStrongPassword123"`
}

// GenericResponse — response umum (sukses tanpa data)
type GenericResponse struct {
	Message string `json:"message" example:"Operation successful"`
}
