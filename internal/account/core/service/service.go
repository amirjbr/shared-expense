package service

import (
	"context"
	"errors"
	"time"

	"github.com/amirjbr/shared-expense/internal/account/app/dto"
	"github.com/amirjbr/shared-expense/internal/account/core/entity"
	"github.com/amirjbr/shared-expense/internal/account/core/port"
	"github.com/amirjbr/shared-expense/internal/account/utils"
	"github.com/amirjbr/shared-expense/pkg/jwt"
	"github.com/amirjbr/shared-expense/pkg/logger"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

//TODO need to work on errors

type UserService struct {
	repo                   port.UserRepo
	logger                 logger.MyLogger
	jwtSecret              []byte
	tokenExpiration        uint
	refreshTokenExpiration uint
}

func NewUserService(repo port.UserRepo, secret []byte, logger logger.MyLogger, tokenExpiration, refreshTokenExpiration uint) (*UserService, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}
	if repo == nil {
		return nil, errors.New("repo is required")
	}
	return &UserService{
		repo:                   repo,
		jwtSecret:              secret,
		logger:                 logger,
		tokenExpiration:        tokenExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}, nil
}

func (s *UserService) Register(ctx context.Context, userReq dto.UserRegisterRequest) (string, error) {

	//TODO add validations and business logics
	//TODO we need to store phone number in 1 way ( for example users input 09011619366 but we store +989011619366)
	//TODO add more rules for password NOT IMPORTANT
	var user entity.User

	user.Username = userReq.Username

	err := utils.ValidatePasswordWithFeedback(userReq.Password)
	if err != nil {
		return "", err
	}
	user.Password = userReq.Password

	err = utils.ValidateEmail(userReq.Email)
	if err != nil {
		return "", errors.New("invalid email")
	}
	user.Email = utils.LowerCaseEmail(userReq.Email)

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

	pass := utils.HashPassword(userReq.Password)
	user.Password = pass

	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.New("user id is required")
	}
	return s.repo.GetUserByID(ctx, id)
}
func (s *UserService) Login(ctx context.Context, loginReq dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	if loginReq.Username == "" {
		return dto.UserLoginResponse{}, errors.New("username is required")
	}
	if loginReq.Password == "" {
		return dto.UserLoginResponse{}, errors.New("password is required")
	}

	user, err := s.repo.GetUserByUsername(ctx, loginReq.Username)
	if err != nil {
		return dto.UserLoginResponse{}, err
	}
	if !utils.CheckPassword(user.Password, loginReq.Password) {
		return dto.UserLoginResponse{}, errors.New("invalid password")
	}

	claims := &jwt.UserClaims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt2.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt2.NewNumericDate(time.Now()),
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(time.Duration(s.tokenExpiration) * time.Minute)),
		},
	}

	token, err := jwt.CreateToken(s.jwtSecret, claims)
	if err != nil {
		return dto.UserLoginResponse{}, err
	}
	return dto.UserLoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	}, nil

}
func (s *UserService) UpdateUser(ctx context.Context, user entity.User) error {
	return s.repo.UpdateUser(ctx, user)
}
