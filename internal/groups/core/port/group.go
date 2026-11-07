package port

import (
	"context"

	"github.com/amirjbr/shared-expense/internal/account/core/service"
	"github.com/amirjbr/shared-expense/internal/groups/core/entity"
	"github.com/google/uuid"
)

type GroupRepo interface {
	CreateGroup(ctx context.Context, group entity.Group) error
	CreateGroupInvitation(ctx context.Context, groupInvitation entity.GroupInvitation) error
	GetGroupByID(ctx context.Context, groupID uuid.UUID) (entity.Group, error)
	IsGroupMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) (bool, error)
	HasPendingInvitation(ctx context.Context, groupID uuid.UUID, invitedUserID uuid.UUID) (bool, error)
	AcceptInvitation(ctx context.Context, InvitationID uuid.UUID, userID uuid.UUID) error
	RejectInvitation(ctx context.Context, InvitationID uuid.UUID, userID uuid.UUID) error
	GetPendingInvitations(ctx context.Context, userID uuid.UUID) ([]entity.GroupInvitation, error)
	GetAllGroupsOfUser(ctx context.Context, userID uuid.UUID) ([]entity.Group, error)
	GetMembersOfGroup(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error)
	AreGroupMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) (bool, error)
}

type UserFinder interface {
	GetUserIDByUsername(ctx context.Context, username string) (uuid.UUID, error)
	GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]service.UserSummary, error)
}
