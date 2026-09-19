package middlewares

import (
	"net/http"

	"github.com/amirjbr/shared-expense/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		claims, err := jwt.ParseToken(token, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"err":     err.Error(),
				"message": "token not valid ",
			})
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
