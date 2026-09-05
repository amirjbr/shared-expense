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

type GroupService struct {
	repo   port.GroupRepo
	logger logger.MyLogger
}

func NewGroupService(repo port.GroupRepo, logger logger.MyLogger) (*GroupService, error) {
	if repo == nil {
		return nil, errors.New("repo is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	return &GroupService{
		repo:   repo,
		logger: logger,
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
