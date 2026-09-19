package routes

import (
	"github.com/amirjbr/shared-expense/internal/account/app/handler"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler *handler.UserHandler
}

func NewUserRoutes(handler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: handler,
	}
}
func (r *UserRoutes) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/register", r.handler.RegisterHandler)
	group.POST("/login", r.handler.LoginHandler)
	group.Use()
}
