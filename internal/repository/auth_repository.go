package repository

import (
	"sort"
	"strings"

	"fooder-backend/internal/dto"
	"fooder-backend/internal/entity"

	"github.com/google/uuid"
)

type AuthRepository interface {
	GetByUsername(username string) (*dto.UserRecord, bool)
	GetAll() (*entity.PaginatedUsers, error)
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

func (r *InMemoryAuthRepository) GetAll() (*entity.PaginatedUsers, error) {
	items := make([]entity.User, 0, len(r.users))
	for _, user := range r.users {
		items = append(items, entity.User{
			ID:       uuid.NewSHA1(uuid.NameSpaceOID, []byte(user.Username)),
			Username: user.Username,
			Password: user.Password,
			Name:     user.Name,
			Roles:    append([]string(nil), user.Roles...),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Username < items[j].Username
	})

	total := len(items)
	totalPages := 0
	if total > 0 {
		totalPages = 1
	}

	return &entity.PaginatedUsers{
		Items:      items,
		TotalItems: total,
		TotalPages: totalPages,
		PageNumber: 1,
		PageSize:   total,
	}, nil
}
