package http

import (
	"github.com/amirjbr/shared-expense/config"
	"github.com/amirjbr/shared-expense/pkg/logger"
	"github.com/gin-gonic/gin"
)

type App struct {
	Server *gin.Engine
	Config config.Config
	Logger logger.MyLogger
}

func NewApp(config config.Config, logger logger.MyLogger) *App {
	var a App
	a.Config = config
	a.Server = gin.Default()
	return &a
}

func (a *App) RunAndListen() {
	// TODO initial handlers here and initial routes then run the server
}
