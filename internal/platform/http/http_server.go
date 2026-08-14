package http

import (
	"log"

	"github.com/amirjbr/shared-expense/config"
	"github.com/amirjbr/shared-expense/internal/account/app/handler"
	"github.com/amirjbr/shared-expense/internal/account/core/service"
	"github.com/amirjbr/shared-expense/pkg/logger"
	"github.com/gin-gonic/gin"
)

type App struct {
	Server  *gin.Engine
	Config  config.Config
	Logger  logger.MyLogger
	UserSvc *service.UserService
}

func NewApp(config config.Config, logger logger.MyLogger, userSvc *service.UserService) *App {
	var a App
	a.Config = config
	a.UserSvc = userSvc
	a.Logger = logger
	a.Server = gin.Default()
	return &a
}

func (a *App) RunAndListen() {

	// TODO initial handlers here and initial routes then run the server
	a.Server = gin.Default()
	h := handler.NewUserHandler(a.UserSvc)
	routes := a.Server.Group("/api/share_expense")
	routes.POST("/register", h.RegisterHandler)

	err := a.Server.Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
		return
	}

}
