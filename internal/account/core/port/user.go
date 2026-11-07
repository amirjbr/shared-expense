package port

import (
	"context"

	"github.com/amirjbr/shared-expense/internal/account/core/entity"
	"github.com/google/uuid"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) error
	GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.User, error)
}
