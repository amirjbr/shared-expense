package service

import (
	"context"
	"errors"

	"github.com/amirjbr/shared-expense/internal/groups/core/port"
	"github.com/google/uuid"
)

type GroupLookup struct {
	repo port.GroupRepo
}

func NewGroupLookup(repo port.GroupRepo) (*GroupLookup, error) {
	if repo == nil {
		return nil, errors.New("nil repo")
	}
	return &GroupLookup{repo: repo}, nil
}
func (s *GroupLookup) AreGroupMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) (bool, error) {

	return s.repo.AreGroupMembers(
		ctx,
		groupID,
		userIDs,
	)
}
