package dto

import "fooder-backend/core/dto"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
}

type UserRecord struct {
	Username string
	Password string
	Name     string
	Roles    []string
}

type PaginatedUsersResponse = dto.Pagination[UserResponse]
