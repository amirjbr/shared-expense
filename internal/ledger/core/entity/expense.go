package entity

import (
	"time"

	"github.com/google/uuid"
)

type Money int64
type SplitType string

const (
	EqualSplitType  SplitType = "equal"
	CustomSplitType SplitType = "custom"
)

type Expense struct {
	ID              uuid.UUID
	GroupID         uuid.UUID
	Title           string
	CreatedByUserID uuid.UUID
	TotalAmount     Money
	SplitType       SplitType
	OccurredAt      time.Time
	Payers          []ExpensePayer
	Shares          []ExpenseShare
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ExpensePayer struct {
	UserID uuid.UUID
	Amount Money
}

type ExpenseShare struct {
	UserID uuid.UUID
	Amount Money
}
