package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func OwnershipGuard(ownerFn func(*gin.Context) (uuid.UUID, error)) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get claims
		val, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "claims missing",
			})
			return
		}

		claims := val.(*Claims)

		ownerID, err := ownerFn(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "failed to resolve owner",
			})
			return
		}

		// 🔥 core check
		if claims.EmployeeID != ownerID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "not allowed",
			})
			return
		}

		c.Next()
	}
}