package test

import (
	"api-shiners/pkg/auth"
	"api-shiners/pkg/entities"
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

//
// ===== MOCK REPOSITORY =====
//
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	user, _ := args.Get(0).(*entities.User)
	return user, args.Error(1)
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) FindRoleByName(ctx context.Context, name string) (*entities.Role, error) {
	args := m.Called(ctx, name)
	role, _ := args.Get(0).(*entities.Role)
	return role, args.Error(1)
}

func (m *MockUserRepo) AssignUserRole(ctx context.Context, ur *entities.UserRole) error {
	args := m.Called(ctx, ur)
	return args.Error(0)
}

func (m *MockUserRepo) SaveResetToken(ctx context.Context, userID uuid.UUID, token string, exp time.Time) error {
	args := m.Called(ctx, userID, token, exp)
	return args.Error(0)
}

func (m *MockUserRepo) FindByResetToken(ctx context.Context, token string) (*entities.User, error) {
	args := m.Called(ctx, token)
	user, _ := args.Get(0).(*entities.User)
	return user, args.Error(1)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// ✅ Tambahan agar test ResetPassword tidak error
func (m *MockUserRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, newHash string) error {
	args := m.Called(ctx, userID, newHash)
	return args.Error(0)
}

func (m *MockUserRepo) ClearResetToken(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

//
// ===== MOCK UTILS =====
//
var sendResetEmailCalled = false

func init() {
	os.Setenv("JWT_SECRET", "testsecret")
	os.Setenv("FRONTEND_URL", "http://localhost:3000")
}

//
// ===== TEST REGISTER =====
//
func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := auth.NewAuthService(mockRepo)

	req := auth.RegisterRequest{
		Name:     "Daffa",
		Email:    "daffa@example.com",
		Password: "123456",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, errors.New("not found"))
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil)
	mockRepo.On("FindRoleByName", mock.Anything, "STUDENT").Return(&entities.Role{ID: uuid.New()}, nil)
	mockRepo.On("AssignUserRole", mock.Anything, mock.AnythingOfType("*entities.UserRole")).Return(nil)

	user, err := service.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Daffa", user.Name)
	assert.Equal(t, "daffa@example.com", user.Email)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := auth.NewAuthService(mockRepo)

	req := auth.RegisterRequest{
		Name:     "Daffa",
		Email:    "daffa@example.com",
		Password: "123456",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(&entities.User{}, nil)

	user, err := service.Register(context.Background(), req)
	assert.Nil(t, user)
	assert.EqualError(t, err, "email already registered")
}

//
// ===== TEST LOGIN =====
//
func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := auth.NewAuthService(mockRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	user := &entities.User{
	ID:           uuid.New(),
	Email:        "daffa@example.com",
	PasswordHash: string(hashed),
	Roles: []*entities.Role{
		{Name: "ADMIN"},
	},
}

	mockRepo.On("FindByEmail", mock.Anything, user.Email).Return(user, nil)

	token, exp, err := service.Login(context.Background(), auth.LoginRequest{
		Email:    user.Email,
		Password: "123456",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, exp.After(time.Now()))
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := auth.NewAuthService(mockRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct123"), bcrypt.DefaultCost)
	user := &entities.User{Email: "daffa@example.com", PasswordHash: string(hashed)}

	mockRepo.On("FindByEmail", mock.Anything, user.Email).Return(user, nil)

	token, exp, err := service.Login(context.Background(), auth.LoginRequest{
		Email:    user.Email,
		Password: "wrongpass",
	})

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, time.Time{}, exp)
}

//
// ===== TEST GENERATE RESET TOKEN =====
//
func TestGenerateResetToken_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := auth.NewAuthService(mockRepo)

	user := &entities.User{
		ID:    uuid.New(),
		Email: "daffa@example.com",
	}

	mockRepo.On("FindByEmail", mock.Anything, user.Email).Return(user, nil)
	mockRepo.On("SaveResetToken", mock.Anything, user.ID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)

	token, err := service.GenerateResetToken(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateResetToken_EmailNotFound(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := auth.NewAuthService(mockRepo)

	mockRepo.On("FindByEmail", mock.Anything, "notfound@example.com").Return(nil, errors.New("not found"))

	token, err := service.GenerateResetToken(context.Background(), "notfound@example.com")
	assert.Error(t, err)
	assert.Empty(t, token)
}

//
// ===== TEST RESET PASSWORD =====
//
func TestResetPassword_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := auth.NewAuthService(mockRepo)

	user := &entities.User{
		ID:           uuid.New(),
		Email:        "daffa@example.com",
		PasswordHash: "oldhash",
	}

	mockRepo.On("FindByResetToken", mock.Anything, "resettoken").Return(user, nil)
	mockRepo.On("UpdatePassword", mock.Anything, user.ID, mock.AnythingOfType("string")).Return(nil)
	mockRepo.On("ClearResetToken", mock.Anything, user.ID).Return(nil)

	err := service.ResetPassword(context.Background(), "resettoken", "newpassword123")
	assert.NoError(t, err)
}

func TestResetPassword_InvalidToken(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := auth.NewAuthService(mockRepo)

	mockRepo.On("FindByResetToken", mock.Anything, "invalidtoken").Return(nil, errors.New("invalid"))

	err := service.ResetPassword(context.Background(), "invalidtoken", "newpass")
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid or expired token")
}

func (m *MockUserRepo) RemoveAllRolesFromUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
