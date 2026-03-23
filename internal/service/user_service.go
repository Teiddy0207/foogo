package service

import (
	"context"
	errors "fooder-backend/core/errors"
	"fooder-backend/core/params"
	"fooder-backend/internal/dto"
	"fooder-backend/internal/repository"
	"fooder-backend/internal/mapper"
)

type UserService struct {
	repo repository.AuthRepository
}

func NewUserService(repo repository.AuthRepository) *UserService {
	return &UserService{repo: repo}
}

func (service *UserService) GetUsers(ctx context.Context, params params.QueryParams) (*dto.PaginatedUsersResponseDto, *errors.AppError) {
	users, err := service.repo.GetAll()
	if err != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch users", err)
	}

	return mapper.ToPaginatedUsersResponse(users), nil
}
