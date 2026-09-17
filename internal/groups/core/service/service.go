package service

import (
	"context"
	"errors"
	"time"

	"github.com/amirjbr/shared-expense/internal/groups/app/dto"
	"github.com/amirjbr/shared-expense/internal/groups/core/entity"
	"github.com/amirjbr/shared-expense/internal/groups/core/port"
	"github.com/amirjbr/shared-expense/pkg/logger"
	"github.com/google/uuid"
)

type CreateInvitationCommand struct {
	Username        string
	GroupID         uuid.UUID
	InvitedByUserID uuid.UUID
}

type GroupService struct {
	repo       port.GroupRepo
	logger     logger.MyLogger
	userFinder port.UserFinder
}

func NewGroupService(repo port.GroupRepo, userFinder port.UserFinder, logger logger.MyLogger) (*GroupService, error) {
	if repo == nil {
		return nil, errors.New("repo is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	return &GroupService{
		repo:       repo,
		logger:     logger,
		userFinder: userFinder,
	}, nil

}

func (s *GroupService) CreateGroup(ctx context.Context, request dto.CreateGroupRequest) (string, error) {
	var group entity.Group

	group.Name = request.Name

	parsedOwnerUUID, err := uuid.Parse(request.OwnerID)
	if err != nil {
		return "", errors.New("invalid owner uuid")
	}
	group.OwnerID = parsedOwnerUUID

	uuidCreated, err := uuid.NewV6()
	if err != nil {
		return "", err
	}
	group.ID = uuidCreated

	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	return s.repo.CreateGroup(ctx, group)
}

func (s *GroupService) InviteMemberToGroup(ctx context.Context, groupInvitationCommand CreateInvitationCommand) error {
	var groupInvitation entity.GroupInvitation
	result, err := s.repo.GetGroupByID(ctx, groupInvitationCommand.GroupID.String())
	if err == nil && result == nil {
		return errors.New("group not found")
	}
	if err != nil {
		return err
	}

	if result.OwnerID != groupInvitationCommand.InvitedByUserID {
		return errors.New("the user is not the owner of the group")
	}

	invitedUserId, err := s.userFinder.GetUserIDByUsername(ctx, groupInvitationCommand.Username)
	if err != nil {
		return err
	}

	isMember, err := s.repo.IsGroupMember(ctx, groupInvitationCommand.GroupID.String(), invitedUserId)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New("user is already a member of the group")
	}

	groupInvitation.ID = uuid.New()
	groupInvitation.CreatedAt = time.Now()

	parseResult, err := uuid.Parse(invitedUserId)
	if err != nil {
		return err
	}

	//TODO for sure need to fix naming of the variables

	groupInvitation.InvitedUserID = parseResult

	groupInvitation.InvitedByUserID = groupInvitationCommand.InvitedByUserID

	groupInvitation.Status = "pending"

	groupInvitation.GroupID = groupInvitationCommand.GroupID

	isHaveInvitation, err := s.repo.IsHaveInvitation(ctx, groupInvitation.InvitedUserID.String(), groupInvitation.GroupID.String())
	if err != nil {
		return err
	}
	if isHaveInvitation {
		return errors.New("the invitation already exist for this group and user")
	}

	return s.repo.CreateGroupInvitation(ctx, groupInvitation)

}

func (s *GroupService) AcceptInvitation(ctx context.Context, invitationID string, userID string) error {
	// TODO complete the service and repo layer also reject in handler and also need refactoring accept handler

	return s.repo.AcceptInvitation(ctx, invitationID, userID)
}

func (s *GroupService) RejectInvitation(ctx context.Context, invitationID string, userID string) error {
	return s.repo.RejectInvitation(ctx, invitationID, userID)
}

func (s *GroupService) GetGroupByID(ctx context.Context, groupID string) (*entity.Group, error) {
	group, err := s.repo.GetGroupByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return group, nil
}
