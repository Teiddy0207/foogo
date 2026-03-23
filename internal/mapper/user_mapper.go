package mapper

import (
	"fooder-backend/internal/dto"
	"fooder-backend/internal/entity"
)

func ToUserDTOResponse(user *dto.UserRecord) *dto.UserResponse {
	response := &dto.UserResponse{
		ID:   0, // This should be set based on your actual user record structure
		Name: user.Name,
		Age:  0, // This should be set based on your actual user record structure
	}

	return response
}

func ToPaginatedUsersResponse(entity *entity.PaginatedUsers) *dto.PaginatedUsersResponseDto {
	if entity == nil {
		return &dto.PaginatedUsersResponseDto{
			Items:      []dto.UserResponse{},
			TotalItems: 0,
			TotalPages: 0,
			PageNumber: 0,
			PageSize:   0,
		}
	}

	items := make([]dto.UserResponse, 0, len(entity.Items))
	for _, user := range entity.Items {
		items = append(items, dto.UserResponse{
			ID:   0,
			Name: user.Name,
			Age:  0,
		})
	}

	return &dto.PaginatedUsersResponseDto{
		Items:      items,
		TotalItems: entity.TotalItems,
		TotalPages: entity.TotalPages,
		PageNumber: entity.PageNumber,
		PageSize:   entity.PageSize,
	}
}
