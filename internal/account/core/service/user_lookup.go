package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

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
	logger   *slog.Logger
}

func NewUserLookUp(userRepo port.UserRepo, logger *slog.Logger) (*UserLookUp, error) {

	if userRepo == nil {
		return nil, errors.New(
			"user repository is nil",
		)
	}

	if logger == nil {
		return nil, errors.New(
			"logger is nil",
		)
	}

	return &UserLookUp{
		userRepo: userRepo,
		logger:   logger,
	}, nil
}

// ---------------------------------------------------------
// Get User ID By Username
// ---------------------------------------------------------

func (s *UserLookUp) GetUserIDByUsername(ctx context.Context, username string) (uuid.UUID, error) {

	username = strings.TrimSpace(username)

	if username == "" {
		return uuid.Nil, errors.New(
			"username is required",
		)
	}

	user, err := s.userRepo.GetUserByUsername(
		ctx,
		username,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"get user by username: %w",
			err,
		)
	}

	return user.ID, nil
}

// ---------------------------------------------------------
// Get Multiple Users
// ---------------------------------------------------------

func (s *UserLookUp) GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]UserSummary, error) {

	if len(ids) == 0 {
		return []UserSummary{}, nil
	}

	for _, id := range ids {
		if id == uuid.Nil {
			return nil, errors.New(
				"user IDs contain an empty UUID",
			)
		}
	}

	users, err := s.userRepo.GetUsersByIDs(
		ctx,
		ids,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get users by IDs: %w",
			err,
		)
	}

	result := make(
		[]UserSummary,
		0,
		len(users),
	)

	for _, user := range users {
		result = append(result,
			UserSummary{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Username:  user.Username,
			},
		)
	}

	return result, nil
}
