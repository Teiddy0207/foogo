package dto

import "fooder-backend/core/dto"

type UserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type PaginatedUsersResponseDto = dto.Pagination[UserResponse]
