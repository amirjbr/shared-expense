package middlewares

import (
	"net/http"
	"strings"

	"github.com/amirjbr/shared-expense/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//func AuthMiddleware(secret []byte) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		token := c.GetHeader("Authorization")
//
//		claims, err := jwt.ParseToken(token, secret)
//		if err != nil {
//			c.AbortWithStatusJSON(
//				http.StatusUnauthorized,
//				gin.H{
//					"message": "invalid authentication token",
//				},
//			)
//			return
//		}
//
//		userID, err := uuid.Parse(claims.UserID)
//		if err != nil {
//			c.AbortWithStatusJSON(
//				http.StatusUnauthorized,
//				gin.H{
//					"message": "invalid authentication token",
//				},
//			)
//			return
//		}
//
//		c.Set("user_id", userID)
//
//		c.Next()
//	}
//}

func AuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "authorization header is required",
				},
			)
			return
		}

		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "invalid authorization header",
				},
			)
			return
		}

		tokenString := strings.TrimPrefix(
			authHeader,
			bearerPrefix,
		)

		tokenString = strings.TrimSpace(tokenString)

		if tokenString == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "authentication token is missing",
				},
			)
			return
		}

		claims, err := jwt.ParseToken(
			tokenString,
			secret,
		)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "invalid authentication token",
				},
			)
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "invalid authentication token",
				},
			)
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
