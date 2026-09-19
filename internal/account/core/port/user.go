package port

import (
	"context"

	"github.com/amirjbr/shared-expense/internal/account/core/entity"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user entity.User) (string, error)
	GetUserByID(ctx context.Context, id string) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) error
}
