package httpcontext

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAuthenticatedUserID(c *gin.Context) (uuid.UUID, bool) {

	value, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}

	userID, ok := value.(uuid.UUID)

	return userID, ok
}
