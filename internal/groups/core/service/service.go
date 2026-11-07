package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	accountservice "github.com/amirjbr/shared-expense/internal/account/core/service"
	"github.com/amirjbr/shared-expense/internal/groups/core/entity"
	"github.com/amirjbr/shared-expense/internal/groups/core/port"
	"github.com/amirjbr/shared-expense/internal/groups/core/service/command"
	"github.com/google/uuid"
)

type GroupService struct {
	repo       port.GroupRepo
	logger     *slog.Logger
	userFinder port.UserFinder
}

func NewGroupService(repo port.GroupRepo, userFinder port.UserFinder, logger *slog.Logger) (*GroupService, error) {

	if repo == nil {
		return nil, errors.New("group repository is nil")
	}

	if userFinder == nil {
		return nil, errors.New("user finder is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	return &GroupService{
		repo:       repo,
		userFinder: userFinder,
		logger:     logger,
	}, nil
}

// ---------------------------------------------------------
// Create Group
// ---------------------------------------------------------

func (s *GroupService) CreateGroup(ctx context.Context, cmd command.CreateGroupCommand) (uuid.UUID, error) {

	if cmd.OwnerID == uuid.Nil {
		return uuid.Nil, errors.New("owner ID is required")
	}

	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return uuid.Nil, errors.New("group name is required")
	}

	now := time.Now().UTC()

	group := entity.Group{
		ID:        uuid.New(),
		Name:      name,
		OwnerID:   cmd.OwnerID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.CreateGroup(ctx, group); err != nil {
		return uuid.Nil, fmt.Errorf(
			"create group: %w",
			err,
		)
	}

	return group.ID, nil
}

// ---------------------------------------------------------
// Invite Member
// ---------------------------------------------------------

func (s *GroupService) InviteMemberToGroup(ctx context.Context, cmd command.CreateInvitationCommand) error {

	if cmd.GroupID == uuid.Nil {
		return errors.New("group ID is required")
	}

	if cmd.InvitedByUserID == uuid.Nil {
		return errors.New("inviting user ID is required")
	}

	username := strings.TrimSpace(cmd.Username)
	if username == "" {
		return errors.New("username is required")
	}

	// Group must exist.
	group, err := s.repo.GetGroupByID(
		ctx,
		cmd.GroupID,
	)
	if err != nil {
		return fmt.Errorf(
			"get group: %w",
			err,
		)
	}

	// Only group owner can invite.
	if group.OwnerID != cmd.InvitedByUserID {
		return errors.New(
			"user is not the owner of the group",
		)
	}

	// Ask Identity module for the user ID.
	invitedUserID, err := s.userFinder.GetUserIDByUsername(
		ctx,
		username,
	)
	if err != nil {
		return fmt.Errorf(
			"find user by username: %w",
			err,
		)
	}

	if invitedUserID == uuid.Nil {
		return errors.New("invalid invited user ID")
	}

	// User must not already belong to the group.
	isMember, err := s.repo.IsGroupMember(
		ctx,
		cmd.GroupID,
		invitedUserID,
	)
	if err != nil {
		return fmt.Errorf(
			"check group membership: %w",
			err,
		)
	}

	if isMember {
		return errors.New(
			"user is already a member of the group",
		)
	}

	// User must not already have a pending invitation.
	hasPendingInvitation, err := s.repo.HasPendingInvitation(
		ctx,
		cmd.GroupID,
		invitedUserID,
	)
	if err != nil {
		return fmt.Errorf(
			"check pending invitation: %w",
			err,
		)
	}

	if hasPendingInvitation {
		return errors.New(
			"a pending invitation already exists for this user and group",
		)
	}

	now := time.Now().UTC()

	invitation := entity.GroupInvitation{
		ID:              uuid.New(),
		GroupID:         cmd.GroupID,
		InvitedUserID:   invitedUserID,
		InvitedByUserID: cmd.InvitedByUserID,
		Status:          "pending",
		CreatedAt:       now,
	}

	if err := s.repo.CreateGroupInvitation(
		ctx,
		invitation,
	); err != nil {
		return fmt.Errorf(
			"create group invitation: %w",
			err,
		)
	}

	return nil
}

// ---------------------------------------------------------
// Accept Invitation
// ---------------------------------------------------------

func (s *GroupService) AcceptInvitation(
	ctx context.Context,
	invitationID uuid.UUID,
	userID uuid.UUID,
) error {

	if invitationID == uuid.Nil {
		return errors.New("invitation ID is required")
	}

	if userID == uuid.Nil {
		return errors.New("user ID is required")
	}

	if err := s.repo.AcceptInvitation(
		ctx,
		invitationID,
		userID,
	); err != nil {
		return fmt.Errorf(
			"accept invitation: %w",
			err,
		)
	}

	return nil
}

// ---------------------------------------------------------
// Reject Invitation
// ---------------------------------------------------------

func (s *GroupService) RejectInvitation(
	ctx context.Context,
	invitationID uuid.UUID,
	userID uuid.UUID,
) error {

	if invitationID == uuid.Nil {
		return errors.New("invitation ID is required")
	}

	if userID == uuid.Nil {
		return errors.New("user ID is required")
	}

	if err := s.repo.RejectInvitation(
		ctx,
		invitationID,
		userID,
	); err != nil {
		return fmt.Errorf(
			"reject invitation: %w",
			err,
		)
	}

	return nil
}

// ---------------------------------------------------------
// Get Group By ID
// ---------------------------------------------------------

func (s *GroupService) GetGroupByID(
	ctx context.Context,
	groupID uuid.UUID,
	userID uuid.UUID,
) (entity.Group, error) {

	if userID == uuid.Nil {
		return entity.Group{}, errors.New(
			"user ID is required",
		)
	}

	if groupID == uuid.Nil {
		return entity.Group{}, errors.New(
			"group ID is required",
		)
	}

	// User must belong to the group.
	isMember, err := s.repo.IsGroupMember(
		ctx,
		groupID,
		userID,
	)
	if err != nil {
		return entity.Group{}, fmt.Errorf(
			"check group membership: %w",
			err,
		)
	}

	if !isMember {
		return entity.Group{}, errors.New(
			"user is not a member of the group",
		)
	}

	group, err := s.repo.GetGroupByID(
		ctx,
		groupID,
	)
	if err != nil {
		return entity.Group{}, fmt.Errorf(
			"get group: %w",
			err,
		)
	}

	return group, nil
}

// ---------------------------------------------------------
// Get Pending Invitations
// ---------------------------------------------------------

func (s *GroupService) GetPendingInvitations(
	ctx context.Context,
	userID uuid.UUID,
) ([]entity.GroupInvitation, error) {

	if userID == uuid.Nil {
		return nil, errors.New(
			"user ID is required",
		)
	}

	invitations, err := s.repo.GetPendingInvitations(
		ctx,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get pending invitations: %w",
			err,
		)
	}

	return invitations, nil
}

// ---------------------------------------------------------
// Get All Groups Of User
// ---------------------------------------------------------

func (s *GroupService) GetAllGroupsOfUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]entity.Group, error) {

	if userID == uuid.Nil {
		return nil, errors.New(
			"user ID is required",
		)
	}

	groups, err := s.repo.GetAllGroupsOfUser(
		ctx,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get groups of user: %w",
			err,
		)
	}

	return groups, nil
}

// ---------------------------------------------------------
// Get Members Of Group
// ---------------------------------------------------------

func (s *GroupService) GetMembersOfGroup(
	ctx context.Context,
	userID uuid.UUID,
	groupID uuid.UUID,
) ([]accountservice.UserSummary, error) {

	if userID == uuid.Nil {
		return nil, errors.New(
			"user ID is required",
		)
	}

	if groupID == uuid.Nil {
		return nil, errors.New(
			"group ID is required",
		)
	}

	// Requesting user must belong to the group.
	isMember, err := s.repo.IsGroupMember(
		ctx,
		groupID,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"check group membership: %w",
			err,
		)
	}

	if !isMember {
		return nil, errors.New(
			"user is not a member of the group",
		)
	}

	userIDs, err := s.repo.GetMembersOfGroup(
		ctx,
		groupID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get group members: %w",
			err,
		)
	}

	users, err := s.userFinder.GetUsersByIDs(
		ctx,
		userIDs,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get member user information: %w",
			err,
		)
	}

	return users, nil
}
