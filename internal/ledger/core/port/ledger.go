package port

import (
	"context"

	"github.com/amirjbr/shared-expense/internal/ledger/core/entity"
	"github.com/google/uuid"
)

type LedgerRepo interface {
	CreateExpense(ctx context.Context, expense entity.Expense) error
}

type GroupFinder interface {
	AreGroupMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) (bool, error)
}
