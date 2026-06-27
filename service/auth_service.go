package service

import (
	"spotsync/dto"
	"spotsync/errors"
	"spotsync/middleware"
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

func (s *authService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	role := req.Role
	if role == "" {
		role = "driver"
	}
	if role == "admin" {
		return nil, errors.ErrForbidden
	}

	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}
	if existingUser != nil {
		return nil, errors.ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.ErrDatabaseError
	}

	return s.mapUserToResponse(user), nil
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}
	if user == nil {
		return nil, errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return &dto.LoginResponse{
		Token: token,
		User:  s.mapUserToLoginResponse(user),
	}, nil
}

func (s *authService) generateJWT(user *models.User) (string, error) {
	claims := &middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
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

func (s *authService) mapUserToLoginResponse(user *models.User) dto.LoginUserResponse {
	return dto.LoginUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}
