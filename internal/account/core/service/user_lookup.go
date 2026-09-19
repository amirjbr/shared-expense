package service

import (
	"context"
	"errors"

	"github.com/amirjbr/shared-expense/internal/account/core/port"
)

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
