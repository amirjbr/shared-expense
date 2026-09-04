package port

import (
	"context"

	"github.com/amirjbr/shared-expense/internal/groups/core/entity"
	"github.com/google/uuid"
)

type GroupRepo interface {
	CreateGroup(ctx context.Context, group entity.Group) (string, error)
	AddMemberToGroup(ctx context.Context, groupID string, userID uuid.UUID) error
}
