package routes

import (
	"github.com/amirjbr/shared-expense/internal/groups/app/handler"
	"github.com/amirjbr/shared-expense/internal/platform/http/middlewares"
	"github.com/amirjbr/shared-expense/pkg/conv"
	"github.com/gin-gonic/gin"
)

type GroupRoutes struct {
	handler *handler.GroupHandler
}

func NewGroupRoutes(handler *handler.GroupHandler) *GroupRoutes {
	return &GroupRoutes{
		handler: handler,
	}
}

func (r *GroupRoutes) RegisterRoutes(group *gin.RouterGroup) {
	group.Use(middlewares.AuthMiddleware(conv.ToBytes("my_secret")))
	group.POST("/groups", r.handler.CreateNewGroupHandler)
}
