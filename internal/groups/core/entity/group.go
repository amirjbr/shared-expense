package entity

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID
	Name      string
	OwnerID   uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
