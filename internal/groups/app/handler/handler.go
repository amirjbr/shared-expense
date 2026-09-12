package handler

import (
	"net/http"

	"github.com/amirjbr/shared-expense/internal/groups/app/dto"
	"github.com/amirjbr/shared-expense/internal/groups/core/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupHandler struct {
	service *service.GroupService
}

func NewGroupHandler(groupSvc *service.GroupService) *GroupHandler {
	return &GroupHandler{
		service: groupSvc,
	}
}

//TODO we can add parse uuid in auth middleware so before it comes to handler it will get checked
//TODO so we dont need to parse everytime in our handlers

func (r *GroupHandler) CreateNewGroupHandler(c *gin.Context) {
	var createGroupReq dto.CreateGroupRequest

	err := c.BindJSON(&createGroupReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createGroupReq.OwnerID = c.Keys["user_id"].(string)

	id, err := r.service.CreateGroup(c, createGroupReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (r *GroupHandler) InviteMemberToGroup(c *gin.Context) {
	var req dto.GroupInvitationRequest
	var groupInvitationCommand service.CreateInvitationCommand
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupInvitationCommand.Username = req.Username

	groupInvitationCommand.InvitedByUserID, err = uuid.Parse(c.Keys["user_id"].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "User ID should be a valid ID",
		})
		return
	}

	groupInvitationCommand.GroupID, err = uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Group ID is not a valid ID",
		})
		return
	}

	err = r.service.InviteMemberToGroup(c, groupInvitationCommand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user invite request sent successfully"})

}

func (r *GroupHandler) AcceptInvitation(c *gin.Context) {
	invitationID, err := uuid.Parse(c.Param("invitation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "invalid invitation id",
		})
		return
	}

	userID, err := uuid.Parse(c.Keys["user_id"].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "invalid user id on jwt",
		})
		return
	}
	err = r.service.AcceptInvitation(c, invitationID.String(), userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user accept invitation request was successfully"})

}

func (r *GroupHandler) RejectInvitation(c *gin.Context) {
	invitationID, err := uuid.Parse(c.Param("invitation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "invalid invitation id",
		})
		return
	}

	userID, err := uuid.Parse(c.Keys["user_id"].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "invalid user id on jwt",
		})
	}

	err = r.service.RejectInvitation(c, invitationID.String(), userID.String())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user reject invitation request was successfully"})

}
