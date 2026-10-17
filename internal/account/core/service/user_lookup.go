package service

import (
	"context"
	"errors"

	"github.com/amirjbr/shared-expense/internal/account/core/port"
	"github.com/google/uuid"
)

type UserSummary struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Username  string
}

type UserLookUp struct {
	userRepo port.UserRepo
}

func NewUserLookUp(userRepo port.UserRepo) (*UserLookUp, error) {
	if userRepo == nil {
		return nil, errors.New("userRepo is nil")
	}
	return &UserLookUp{userRepo: userRepo}, nil
}

func (s *UserLookUp) GetUserIDByUsername(ctx context.Context, username string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	return user.ID.String(), nil
}

func (s *UserLookUp) GetUsersByIDs(ctx context.Context, ids []string) ([]UserSummary, error) {
	if len(ids) == 0 {
		return []UserSummary{}, nil
	}
	var users []UserSummary
	userRepoResult, err := s.userRepo.GetUsersByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, user := range userRepoResult {
		var u UserSummary
		u.ID = user.ID
		u.FirstName = user.FirstName
		u.LastName = user.LastName
		u.Username = user.Username
		users = append(users, u)
	}

	return users, nil

}
