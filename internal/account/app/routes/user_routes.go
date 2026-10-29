package routes

import (
	"log/slog"

	"github.com/amirjbr/shared-expense/internal/account/app/handler"
	"github.com/amirjbr/shared-expense/internal/platform/http/middlewares"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler *handler.UserHandler
	logger  *slog.Logger
}

func NewUserRoutes(handler *handler.UserHandler, logger *slog.Logger) *UserRoutes {
	return &UserRoutes{
		handler: handler,
		logger:  logger,
	}
}
func (r *UserRoutes) RegisterRoutes(group *gin.RouterGroup) {
	group.Use(middlewares.RequestIDMiddleware())
	group.Use(middlewares.RequestLoggerMiddleware(r.logger))
	group.POST("/register", r.handler.RegisterHandler)
	group.POST("/login", r.handler.LoginHandler)
	group.Use()
}
