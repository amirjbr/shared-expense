package handler

import "github.com/amirjbr/shared-expense/internal/ledger/core/service"

type LedgerHandler struct {
	service *service.LedgerService
}

func NewLedgerHandler(service *service.LedgerService) *LedgerHandler {
	return &LedgerHandler{
		service: service,
	}
}
