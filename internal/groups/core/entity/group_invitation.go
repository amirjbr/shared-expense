package entity

import (
	"time"

	"github.com/google/uuid"
)

type GroupInvitation struct {
	ID              uuid.UUID
	GroupID         uuid.UUID
	InvitedUserID   uuid.UUID
	InvitedByUserID uuid.UUID
	Status          GroupInvitationStatus
	CreatedAt       time.Time
	RespondedAt     time.Time
}

type GroupInvitationStatus string

const (
	GroupInvitationStatusPending  GroupInvitationStatus = "PENDING"
	GroupInvitationStatusAccepted GroupInvitationStatus = "ACCEPTED"
	GroupInvitationStatusRejected GroupInvitationStatus = "REJECTED"
)
