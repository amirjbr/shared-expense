package handler

import (
	"net/http"

	"github.com/amirjbr/shared-expense/internal/account/app/dto"
	"github.com/amirjbr/shared-expense/internal/account/core/service"
	"github.com/amirjbr/shared-expense/internal/account/core/service/command"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {

	return &UserHandler{
		userSvc: userSvc,
	}
}

// ---------------------------------------------------------
// Register
// ---------------------------------------------------------

func (h *UserHandler) RegisterHandler(c *gin.Context) {
	var req dto.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid request body",
			},
		)
		return
	}

	cmd := command.RegisterUserCommand{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Username:    req.Username,
	}

	userID, err := h.userSvc.Register(c.Request.Context(), cmd)
	if err != nil {
		/*
			Later your centralized error handling will
			distinguish validation/conflict/database errors.
		*/
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to register user",
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"id": userID.String(),
		},
	)
}

// ---------------------------------------------------------
// Login
// ---------------------------------------------------------

func (h *UserHandler) LoginHandler(c *gin.Context) {
	var req dto.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid request body",
			},
		)
		return
	}

	cmd := command.LoginUserCommand{
		Username: req.Username,
		Password: req.Password,
	}

	result, err := h.userSvc.Login(c.Request.Context(), cmd)
	if err != nil {
		/*
			When proper error mapping is added,
			invalid credentials should become 401,
			not 500.
		*/
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to login",
			},
		)
		return
	}

	response := dto.UserLoginResponse{
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}
