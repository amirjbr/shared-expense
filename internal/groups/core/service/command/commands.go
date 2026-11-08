package command

import "github.com/google/uuid"

type CreateGroupCommand struct {
	Name    string
	OwnerID uuid.UUID
}

type CreateInvitationCommand struct {
	Username        string
	InvitedByUserID uuid.UUID
	GroupID         uuid.UUID
}
