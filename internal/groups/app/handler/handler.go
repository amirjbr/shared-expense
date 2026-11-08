package handler

import (
	"net/http"

	"github.com/amirjbr/shared-expense/internal/groups/app/dto"
	"github.com/amirjbr/shared-expense/internal/groups/core/service"
	"github.com/amirjbr/shared-expense/internal/groups/core/service/command"
	"github.com/amirjbr/shared-expense/internal/platform/http/httpcontext"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupHandler struct {
	service *service.GroupService
}

func NewGroupHandler(
	groupSvc *service.GroupService,
) *GroupHandler {
	return &GroupHandler{
		service: groupSvc,
	}
}

// ---------------------------------------------------------
// Create Group
// ---------------------------------------------------------

func (r *GroupHandler) CreateNewGroupHandler(c *gin.Context) {
	var req dto.CreateGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	ownerID, ok := httpcontext.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "authenticated user information is unavailable",
		})
		return
	}

	cmd := command.CreateGroupCommand{
		Name:    req.Name,
		OwnerID: ownerID,
	}

	groupID, err := r.service.CreateGroup(
		c.Request.Context(),
		cmd,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create group",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": groupID.String(),
	})
}

// ---------------------------------------------------------
// Invite Member
// ---------------------------------------------------------

func (r *GroupHandler) InviteMemberToGroup(c *gin.Context) {
	var req dto.GroupInvitationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid group ID",
		})
		return
	}

	invitedByUserID, ok := httpcontext.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "authenticated user information is unavailable",
		})
		return
	}

	cmd := command.CreateInvitationCommand{
		Username:        req.Username,
		InvitedByUserID: invitedByUserID,
		GroupID:         groupID,
	}

	if err := r.service.InviteMemberToGroup(
		c.Request.Context(),
		cmd,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to send group invitation",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "group invitation created successfully",
	})
}

// ---------------------------------------------------------
// Accept Invitation
// ---------------------------------------------------------

func (r *GroupHandler) AcceptInvitation(c *gin.Context) {
	invitationID, err := uuid.Parse(
		c.Param("invitation_id"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid invitation ID",
		})
		return
	}

	userID, ok := httpcontext.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "authenticated user information is unavailable",
		})
		return
	}

	if err := r.service.AcceptInvitation(
		c.Request.Context(),
		invitationID,
		userID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to accept invitation",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "invitation accepted successfully",
	})
}

// ---------------------------------------------------------
// Reject Invitation
// ---------------------------------------------------------

func (r *GroupHandler) RejectInvitation(c *gin.Context) {
	invitationID, err := uuid.Parse(
		c.Param("invitation_id"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid invitation ID",
		})
		return
	}

	userID, ok := httpcontext.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "authenticated user information is unavailable",
		})
		return
	}

	if err := r.service.RejectInvitation(
		c.Request.Context(),
		invitationID,
		userID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to reject invitation",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "invitation rejected successfully",
	})
}

// ---------------------------------------------------------
// Get Group Information
// ---------------------------------------------------------

func (r *GroupHandler) GetGroupInformation(c *gin.Context) {
	groupID, err := uuid.Parse(
		c.Param("group_id"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid group ID",
		})
		return
	}

	userID, ok := httpcontext.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "authenticated user information is unavailable",
		})
		return
	}

	group, err := r.service.GetGroupByID(
		c.Request.Context(),
		groupID,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get group information",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"group": group,
	})
}

// ---------------------------------------------------------
// Get Pending Invitations
// ---------------------------------------------------------

func (r *GroupHandler) GetPendingInvitationOfUser(
	c *gin.Context,
) {
	userID, ok := httpcontext.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "authenticated user information is unavailable",
		})
		return
	}

	invitations, err := r.service.GetPendingInvitations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get pending invitations",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"invitations": invitations,
	})
}

// ---------------------------------------------------------
// Get Groups Of Authenticated User
// ---------------------------------------------------------

func (r *GroupHandler) GetGroupsOfAuthenticatedUser(
	c *gin.Context,
) {
	userID, ok := httpcontext.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "authenticated user information is unavailable",
		})
		return
	}

	groups, err := r.service.GetAllGroupsOfUser(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get groups",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"groups": groups,
	})
}

// ---------------------------------------------------------
// Get Members Of Group
// ---------------------------------------------------------

func (r *GroupHandler) GetMembersOfGroup(c *gin.Context) {
	groupID, err := uuid.Parse(
		c.Param("group_id"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid group ID",
		})
		return
	}

	userID, ok := httpcontext.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "authenticated user information is unavailable",
		})
		return
	}

	members, err := r.service.GetMembersOfGroup(
		c.Request.Context(),
		userID,
		groupID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get group members",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"members": members,
	})
}
