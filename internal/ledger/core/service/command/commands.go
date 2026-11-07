package command

import (
	"time"

	"github.com/amirjbr/shared-expense/internal/ledger/core/entity"
	"github.com/google/uuid"
)

type CreateExpenseCommand struct {
	GroupID         uuid.UUID
	CreatedByUserID uuid.UUID

	Title       string
	TotalAmount entity.Money
	SplitType   entity.SplitType
	OccurredAt  time.Time

	Payers       []CreateExpensePayer
	Participants []uuid.UUID
	Shares       []CreateExpenseShare
}

type CreateExpensePayer struct {
	UserID uuid.UUID
	Amount entity.Money
}
type CreateExpenseShare struct {
	UserID uuid.UUID
	Amount entity.Money
}
