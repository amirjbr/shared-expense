package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/amirjbr/shared-expense/internal/account/app/dto"
	"github.com/amirjbr/shared-expense/internal/account/core/entity"
	"github.com/amirjbr/shared-expense/internal/account/core/port"
	"github.com/amirjbr/shared-expense/pkg/conv"
	"github.com/amirjbr/shared-expense/pkg/logger"
	"github.com/google/uuid"
)

type UserService struct {
	repo   port.UserRepo
	logger logger.MyLogger
}

func NewUserService(repo port.UserRepo, logger logger.MyLogger) (*UserService, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}
	if repo == nil {
		return nil, errors.New("repo is required")
	}
	return &UserService{
		repo:   repo,
		logger: logger,
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, userReq dto.UserRegisterRequest) (string, error) {

	//TODO add validations and business logics
	var user entity.User

	user.Username = userReq.Username
	user.Password = userReq.Password
	user.Email = userReq.Email
	user.PhoneNumber = userReq.PhoneNumber
	user.FirstName = userReq.FirstName
	user.LastName = userReq.LastName
	uuidCreated, err := uuid.NewV6()
	if err != nil {
		return "", err
	}
	user.ID = uuidCreated

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	pass := HashPassword(userReq.Password)
	user.Password = pass

	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.New("user id is required")
	}
	return s.repo.GetUserByID(ctx, id)
}
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	return s.repo.GetUserByUsername(ctx, username)
}
func (s *UserService) UpdateUser(ctx context.Context, user entity.User) error {
	return s.repo.UpdateUser(ctx, user)
}
func HashPassword(pass string) string {
	h := sha256.New()
	h.Write(conv.ToBytes(pass))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
