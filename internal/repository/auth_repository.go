package repository

import (
	"strings"

	"fooder-backend/internal/dto"
)

type AuthRepository interface {
	GetByUsername(username string) (*dto.UserRecord, bool)
}

type InMemoryAuthRepository struct {
	users map[string]dto.UserRecord
}

func NewInMemoryAuthRepository() *InMemoryAuthRepository {
	return &InMemoryAuthRepository{
		users: map[string]dto.UserRecord{
			"admin": {
				Username: "admin",
				Password: "admin123",
				Name:     "Administrator",
				Roles:    []string{"admin"},
			},
			"user": {
				Username: "user",
				Password: "user123",
				Name:     "Normal User",
				Roles:    []string{"user"},
			},
		},
	}
}

func (r *InMemoryAuthRepository) GetByUsername(username string) (*dto.UserRecord, bool) {
	key := strings.TrimSpace(strings.ToLower(username))
	if key == "" {
		return nil, false
	}

	user, ok := r.users[key]
	if !ok {
		return nil, false
	}

	copyUser := user
	return &copyUser, true
}
