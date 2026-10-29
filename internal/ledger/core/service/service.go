package service

import (
	"errors"
	"log/slog"

	"github.com/amirjbr/shared-expense/internal/ledger/core/port"
)

type LedgerService struct {
	repo   *port.LedgerRepo
	logger *slog.Logger
}

func NewLedgerService(repo *port.LedgerRepo, logger *slog.Logger) (*LedgerService, error) {
	if repo == nil {
		return nil, errors.New("repo is nil")
	}
	return &LedgerService{repo: repo, logger: logger}, nil
}
