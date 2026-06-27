package service

import (
	"spotsync/dto"
	"spotsync/errors"
	"spotsync/models"
	"spotsync/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user with bcrypt hashed password
func (s *authService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}
	if existingUser != nil {
		return nil, errors.ErrUserExists
	}

	// Hash password with bcrypt (cost 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.ErrDatabaseError
	}

	return s.mapUserToResponse(user), nil
}

// Login validates credentials and returns JWT token
func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}
	if user == nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return &dto.LoginResponse{
		Token: token,
		User:  *s.mapUserToResponse(user),
	}, nil
}

// generateJWT creates a signed JWT token with user claims
func (s *authService) generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours expiry
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) mapUserToResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
