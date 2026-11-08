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
	group.POST("/groups/:group_id/expenses", r.handler.CreateExpenseHandler)
	//group.GET("/groups/:group_id/expenses", r.handler.GetGroupExpenseHistory)
	//group.GET("/groups/:group_id/expenses/:expense_id", r.handler.GetAnExpenseHistoryOfGroup)
	//group.GET("/groups/:group_id/balance", r.handler.CalculateCurrentMemberBalancesOfGroup)
	//group.POST("/groups/:group_id/settlement", r.handler.CreateSettlementRecord)
	//group.GET("/groups/:group_id/settlements", r.handler.GetAllSettlementsRecords)
}
