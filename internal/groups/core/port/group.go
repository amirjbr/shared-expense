package port

import (
	"context"

	"github.com/amirjbr/shared-expense/internal/groups/core/entity"
)

type GroupRepo interface {
	CreateGroup(ctx context.Context, group entity.Group) (string, error)
	CreateGroupInvitation(ctx context.Context, groupInvitation entity.GroupInvitation) error
	GetGroupByID(ctx context.Context, groupID string) (*entity.Group, error)
	IsGroupMember(ctx context.Context, groupID string, userID string) (bool, error)
	IsHaveInvitation(ctx context.Context, invitedUserID string, groupID string) (bool, error)
	AcceptInvitation(ctx context.Context, InvitationID string, userID string) error
	RejectInvitation(ctx context.Context, InvitationID string, userID string) error
}

type UserFinder interface {
	GetUserIDByUsername(ctx context.Context, username string) (string, error)
}
