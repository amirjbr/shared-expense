package dto

type CreateGroupRequest struct {
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type GroupInvitationRequest struct {
	Username string `json:"username""`
}
