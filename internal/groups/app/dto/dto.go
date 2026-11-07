package dto

type CreateGroupRequest struct {
	Name string `json:"name"`
}

type GroupInvitationRequest struct {
	Username string `json:"username"`
}
