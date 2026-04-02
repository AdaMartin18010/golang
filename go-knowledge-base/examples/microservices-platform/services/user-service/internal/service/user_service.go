package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"user-service/internal/domain"
)

type UserService interface {
	Register(req *domain.CreateUserRequest) (*domain.AuthResponse, error)
	Login(req *domain.LoginRequest) (*domain.AuthResponse, error)
	GetProfile(userID uuid.UUID) (*domain.User, error)
	UpdateProfile(userID uuid.UUID, req *domain.UpdateUserRequest) (*domain.User, error)
	RefreshToken(refreshToken string) (*domain.AuthResponse, error)
	ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error
	DeleteAccount(userID uuid.UUID) error
	VerifyEmail(token string) error
}

type userService struct {
	repo       domain.UserRepository
	jwtSecret  string
	jwtExpiry  time.Duration
}

func NewUserService(repo domain.UserRepository, jwtSecret string, jwtExpiry time.Duration) UserService {
	return &userService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func (s *userService) Register(req *domain.CreateUserRequest) (*domain.AuthResponse, error) {
	// Check if email exists
	exists, err := s.repo.Exists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailExists
	}

	// Create user
	user := &domain.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtExpiry.Seconds()),
		User:         user,
	}, nil
}

func (s *userService) Login(req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.CheckPassword(req.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtExpiry.Seconds()),
		User:         user,
	}, nil
}

func (s *userService) GetProfile(userID uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(userID)
}

func (s *userService) UpdateProfile(userID uuid.UUID, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) RefreshToken(refreshToken string) (*domain.AuthResponse, error) {
	// Parse and validate refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret + "_refresh"), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	user, err := s.repo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	newAccessToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.jwtExpiry.Seconds()),
		User:         user,
	}, nil
}

func (s *userService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}

	if !user.CheckPassword(oldPassword) {
		return errors.New("invalid old password")
	}

	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	// Update password hash - would need a separate method in repo
	return nil
}

func (s *userService) DeleteAccount(userID uuid.UUID) error {
	return s.repo.Delete(userID)
}

func (s *userService) VerifyEmail(token string) error {
	// Implementation for email verification
	return nil
}

func (s *userService) generateToken(user *domain.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *userService) generateRefreshToken(user *domain.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret + "_refresh"))
}
