package entity

import (
	"fooder-backend/core/entity"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
	Name     string    `db:"name"`
	Roles    []string  `db:"roles"`
}

type PaginatedUsers = entity.Pagination[User]
