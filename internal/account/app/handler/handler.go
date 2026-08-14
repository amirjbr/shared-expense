package handler

import (
	"net/http"

	"github.com/amirjbr/shared-expense/internal/account/app/dto"
	"github.com/amirjbr/shared-expense/internal/account/core/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}

func (h *Handler) RegisterHandler(c *gin.Context) {
	var registerReq dto.UserRegisterRequest
	err := c.BindJSON(&registerReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	createdId, err := h.userSvc.CreateUser(c, registerReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"id": createdId})
}
func (h *Handler) LoginHandler(c *gin.Context) {
	var loginReq dto.UserLoginRequest
	err := c.BindJSON(&loginReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := h.userSvc.GetUserByUsername(c, loginReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, user)
}
