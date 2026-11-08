package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/amirjbr/shared-expense/internal/account/core/entity"
	"github.com/amirjbr/shared-expense/internal/account/core/port"
	"github.com/amirjbr/shared-expense/internal/account/core/service/command"
	"github.com/amirjbr/shared-expense/internal/account/utils"
	"github.com/amirjbr/shared-expense/pkg/jwt"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserService struct {
	repo                   port.UserRepo
	logger                 *slog.Logger
	jwtSecret              []byte
	tokenExpiration        uint
	refreshTokenExpiration uint
}

type LoginResult struct {
	AccessToken string
	TokenType   string
}

func NewUserService(repo port.UserRepo, secret []byte, logger *slog.Logger,
	tokenExpiration uint, refreshTokenExpiration uint) (*UserService, error) {

	if repo == nil {
		return nil, errors.New("user repository is required")
	}

	if logger == nil {
		return nil, errors.New("logger is required")
	}

	if len(secret) == 0 {
		return nil, errors.New("JWT secret is required")
	}

	return &UserService{
		repo:                   repo,
		jwtSecret:              secret,
		logger:                 logger,
		tokenExpiration:        tokenExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}, nil
}

// ---------------------------------------------------------
// Register
// ---------------------------------------------------------

func (s *UserService) Register(ctx context.Context, cmd command.RegisterUserCommand) (uuid.UUID, error) {

	username := strings.TrimSpace(cmd.Username)
	if username == "" {
		return uuid.Nil, errors.New("username is required")
	}

	firstName := strings.TrimSpace(cmd.FirstName)
	if firstName == "" {
		return uuid.Nil, errors.New("first name is required")
	}

	lastName := strings.TrimSpace(cmd.LastName)
	if lastName == "" {
		return uuid.Nil, errors.New("last name is required")
	}

	email := strings.TrimSpace(cmd.Email)
	if email == "" {
		return uuid.Nil, errors.New("email is required")
	}

	if strings.TrimSpace(cmd.Password) == "" {
		return uuid.Nil, errors.New("password is required")
	}

	if err := utils.ValidatePasswordWithFeedback(
		cmd.Password,
	); err != nil {
		return uuid.Nil, err
	}

	if err := utils.ValidateEmail(email); err != nil {
		return uuid.Nil, errors.New("invalid email")
	}

	normalizedEmail := utils.LowerCaseEmail(email)

	/*
		TODO:
		Normalize phone number into one canonical format,
		for example:
		09011619366 → +989011619366
	*/

	now := time.Now().UTC()

	user := entity.User{
		ID:          uuid.New(),
		Username:    username,
		Password:    utils.HashPassword(cmd.Password),
		Email:       normalizedEmail,
		PhoneNumber: strings.TrimSpace(cmd.PhoneNumber),
		FirstName:   firstName,
		LastName:    lastName,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return uuid.Nil, fmt.Errorf(
			"create user: %w",
			err,
		)
	}

	return user.ID, nil
}

// ---------------------------------------------------------
// Login
// ---------------------------------------------------------

func (s *UserService) Login(ctx context.Context, cmd command.LoginUserCommand) (LoginResult, error) {

	username := strings.TrimSpace(cmd.Username)

	if username == "" {
		return LoginResult{}, errors.New(
			"username is required",
		)
	}

	if cmd.Password == "" {
		return LoginResult{}, errors.New(
			"password is required",
		)
	}

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return LoginResult{}, fmt.Errorf(
			"get user by username: %w",
			err,
		)
	}

	if !utils.CheckPassword(user.Password, cmd.Password) {
		return LoginResult{}, errors.New(
			"invalid credentials",
		)
	}

	now := time.Now().UTC()

	claims := &jwt.UserClaims{
		UserID: user.ID.String(),

		RegisteredClaims: jwt2.RegisteredClaims{
			Subject: user.ID.String(),

			IssuedAt: jwt2.NewNumericDate(
				now,
			),

			ExpiresAt: jwt2.NewNumericDate(
				now.Add(
					time.Duration(
						s.tokenExpiration,
					) * time.Minute,
				),
			),
		},
	}

	token, err := jwt.CreateToken(s.jwtSecret, claims)
	if err != nil {
		return LoginResult{}, fmt.Errorf(
			"create access token: %w",
			err,
		)
	}

	return LoginResult{
		AccessToken: token,
		TokenType:   "Bearer",
	}, nil
}

// ---------------------------------------------------------
// Update User
// ---------------------------------------------------------

func (s *UserService) UpdateUser(ctx context.Context, user entity.User) error {

	if user.ID == uuid.Nil {
		return errors.New("user ID is required")
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf(
			"update user: %w",
			err,
		)
	}

	return nil
}
