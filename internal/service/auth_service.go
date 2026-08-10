package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"example.com/my-api/config"
	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/model"
	"example.com/my-api/internal/repository"
)

var ErrInvalidCredentials = errors.New("email atau password salah")
var ErrEmailExists = errors.New("email sudah terdaftar")

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (string, string, error)
	Login(ctx context.Context, req dto.LoginRequest) (string, error)
}

type authService struct {
	repository repository.UserRepository
	config     config.Config
}

func NewAuthService(repository repository.UserRepository, cfg config.Config) AuthService {
	return &authService{repository: repository, config: cfg}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (string, string, error) {
	_, err := s.repository.FindByEmail(ctx, req.Email)
	if err == nil {
		return "", "", ErrEmailExists
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return "", "", err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	user, err := s.repository.Create(ctx, model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	})
	if err != nil {
		if errors.Is(err, repository.ErrEmailExists) {
			return "", "", ErrEmailExists
		}
		return "", "", err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", "", err
	}

	return token, user.Email, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	user, err := s.repository.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) generateToken(user model.User) (string, error) {
	expiry, err := time.ParseDuration(s.config.JWTExpiry)
	if err != nil || expiry <= 0 {
		expiry = 24 * time.Hour
	}

	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}
