package auth

import (
	"api-shiners/pkg/entities"
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"api-shiners/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*entities.User, error)
	Login(ctx context.Context, req LoginRequest) (string, error)
	Logout(ctx context.Context, token string) error
	GenerateResetToken(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type authService struct {
	userRepo UserRepository
}

func NewAuthService(userRepo UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// =============================
// REGISTER
// =============================
func (s *authService) Register(ctx context.Context, req RegisterRequest) (*entities.User, error) {
	// Cek duplikasi email
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := &entities.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashed),
		IsActive:     true,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Assign role default: STUDENT
	role, err := s.userRepo.FindRoleByName(ctx, string(entities.STUDENT))
	if err == nil {
		_ = s.userRepo.AssignUserRole(ctx, &entities.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		})
	}

	return user, nil
}

// =============================
// LOGIN (Generate JWT Token)
// =============================
func (s *authService) Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Ambil secret key dari env
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set in environment")
	}

	// === Remember Me Logic ===
	// Default durasi token: 24 jam
	expiration := time.Now().Add(24 * time.Hour)

	// Jika JWT_EXPIRE_HOURS diset di .env
	if os.Getenv("JWT_EXPIRE_HOURS") != "" {
		if d, err := time.ParseDuration(os.Getenv("JWT_EXPIRE_HOURS") + "h"); err == nil {
			expiration = time.Now().Add(d)
		}
	}

	// Jika request menyertakan remember_me = true → token berlaku 7 hari
	rememberMe := false
	if val, ok := ctx.Value("remember_me").(bool); ok {
		rememberMe = val
	}
	if rememberMe {
		expiration = time.Now().Add(7 * 24 * time.Hour) // 7 hari
	}

	// Buat payload token
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     expiration.Unix(),
		"iat":     time.Now().Unix(),
		"roles":   user.Roles,
	}

	// Generate token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return signedToken, nil
}


// LOGOUT
func (s *authService) Logout(ctx context.Context, token string) error {
	return nil
}

// =============================
// RESET PASSWORD
// =============================
func (s *authService) GenerateResetToken(ctx context.Context, email string) (string, error) {
	// Cek apakah email terdaftar
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("email not found")
	}

	// Buat token baru dan waktu kadaluarsa (1 jam)
	token := uuid.NewString()
	expiresAt := time.Now().Add(1 * time.Hour)

	// Simpan token ke database
	if err := s.userRepo.SaveResetToken(ctx, user.ID, token, expiresAt); err != nil {
		return "", fmt.Errorf("failed to save reset token: %v", err)
	}

	// Kirim email ke user
	if err := utils.SendResetEmail(user.Email, token); err != nil {
		return "", fmt.Errorf("failed to send reset email: %v", err)
	}

	fmt.Printf("📨 Reset password email sent to %s\n", user.Email)
	fmt.Printf("🔗 Reset link: %s/reset-password?token=%s\n", os.Getenv("FRONTEND_URL"), token)

	return token, nil
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Temukan user berdasarkan token (dan pastikan belum expired)
	user, err := s.userRepo.FindByResetToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// Hash password baru
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, user.ID, string(hashed)); err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	// Hapus token setelah digunakan
	if err := s.userRepo.ClearResetToken(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to clear reset token: %v", err)
	}

	fmt.Printf("✅ Password successfully reset for user %s\n", user.Email)
	return nil
}

