package entity

import (
	"time"

	"github.com/google/uuid"
)

type Settlement struct {
	ID              uuid.UUID
	GroupID         uuid.UUID
	PayerUserID     uuid.UUID
	ReceiverID      uuid.UUID
	Amount          Money
	SettledAt       time.Time
	CreatedByUserID uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
