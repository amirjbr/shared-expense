package handler

import (
	"net/http"

	"github.com/amirjbr/shared-expense/internal/groups/app/dto"
	"github.com/amirjbr/shared-expense/internal/groups/core/service"
	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	service *service.GroupService
}

func NewGroupHandler(groupSvc *service.GroupService) *GroupHandler {
	return &GroupHandler{
		service: groupSvc,
	}
}

func (r *GroupHandler) CreateNewGroupHandler(c *gin.Context) {
	var createGroupReq dto.CreateGroupRequestWithoutID

	err := c.BindJSON(&createGroupReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var createGroupRequestWithID dto.CreateGroupRequest
	createGroupRequestWithID.Name = createGroupReq.Name
	createGroupRequestWithID.OwnerID = c.Keys["user_id"].(string)

	id, err := r.service.CreateGroup(c, createGroupRequestWithID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
