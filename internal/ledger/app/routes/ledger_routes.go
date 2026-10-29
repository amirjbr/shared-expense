package routes

import (
	"log/slog"

	"github.com/amirjbr/shared-expense/internal/ledger/app/handler"
	"github.com/amirjbr/shared-expense/internal/platform/http/middlewares"
	"github.com/amirjbr/shared-expense/pkg/conv"
	"github.com/gin-gonic/gin"
)

type LedgerRoutes struct {
	handler *handler.LedgerHandler
	logger  *slog.Logger
}

func NewLedgerRoutes(handler *handler.LedgerHandler, logger *slog.Logger) *LedgerRoutes {
	return &LedgerRoutes{
		handler: handler,
		logger:  logger,
	}
}

func (r *LedgerRoutes) RegisterRoutes(group *gin.RouterGroup) {
	group.Use(middlewares.RequestIDMiddleware())
	group.Use(middlewares.RequestLoggerMiddleware(r.logger))
	group.Use(middlewares.AuthMiddleware(conv.ToBytes("my_secret")))
}
