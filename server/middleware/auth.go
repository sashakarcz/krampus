package middleware

import (
	"krampus/server/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and injects user info into context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Try to extract token from Authorization header first
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Check Bearer token format
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// If not in header, try to get from cookie
		if tokenString == "" {
			cookie, err := c.Cookie("token")
			if err == nil && cookie != "" {
				tokenString = cookie
			}
		}

		// If still no token found, return unauthorized
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := services.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Inject user info into context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(int64)
	return id, ok
}

// GetUsername retrieves username from context
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}

// GetRole retrieves user role from context
func GetRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	r, ok := role.(string)
	return r, ok
}
