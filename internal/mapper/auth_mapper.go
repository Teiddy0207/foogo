package mapper

import "fooder-backend/internal/dto"

func ToUserResponse(user *dto.UserRecord) *dto.UserResponse {
	return &dto.UserResponse{
		Username: user.Username,
		Name:     user.Name,
		Roles:    user.Roles,
	}
}

func ToLoginResponse(user *dto.UserRecord, tokenPrefix string) *dto.LoginResponse {
	return &dto.LoginResponse{
		Token: tokenPrefix + user.Username,
		User:  *ToUserResponse(user),
	}
}
