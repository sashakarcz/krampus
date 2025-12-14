package middleware

import (
	"krampus/server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware checks if the user has admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := GetRole(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		if role != string(models.RoleAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
