package http

import (
	"github.com/amirjbr/shared-expense/config"
	"github.com/gin-gonic/gin"
)

type App struct {
	Server *gin.Engine
	Config config.Config
}

func NewApp(config config.Config) *App {
	var a App
	a.Config = config
	a.Server = gin.Default()
	return &a
}

func (a *App) RunAndListen() {
	// TODO initial handlers here and initial routes then run the server
}
