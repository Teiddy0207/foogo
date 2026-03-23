package service

import (
	"crypto/subtle"
	"strings"

	errors "fooder-backend/core/errors"
	"fooder-backend/internal/dto"
	"fooder-backend/internal/mapper"
	"fooder-backend/internal/repository"
)

type AuthService struct {
	repo           repository.AuthRepository
	devTokenPrefix string
}

func NewAuthService(repo repository.AuthRepository, devTokenPrefix string) *AuthService {
	return &AuthService{repo: repo, devTokenPrefix: devTokenPrefix}
}

func (s *AuthService) Login(input dto.LoginRequest) (*dto.LoginResponse, *errors.AppError) {
	username := strings.TrimSpace(strings.ToLower(input.Username))
	password := strings.TrimSpace(input.Password)

	if username == "" || password == "" {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "username and password are required", nil)
	}

	user, ok := s.repo.GetByUsername(username)
	if !ok {
		return nil, errors.NewAppError(errors.ErrInvalidCredentials, "invalid credentials", nil)
	}

	if subtle.ConstantTimeCompare([]byte(user.Password), []byte(password)) != 1 {
		return nil, errors.NewAppError(errors.ErrInvalidCredentials, "invalid credentials", nil)
	}

	return mapper.ToLoginResponse(user, s.devTokenPrefix), nil
}

func (s *AuthService) Me(token string) (*dto.AuthUserResponse, *errors.AppError) {
	t := strings.TrimSpace(token)
	prefix := s.devTokenPrefix
	if prefix == "" || !strings.HasPrefix(t, prefix) {
		return nil, errors.NewAppError(errors.ErrInvalidTokenFormat, "invalid token", nil)
	}

	username := strings.TrimPrefix(t, prefix)
	user, ok := s.repo.GetByUsername(username)
	if !ok {
		return nil, errors.NewAppError(errors.ErrUnauthorized, "invalid token", nil)
	}

	return mapper.ToUserResponse(user), nil
}
