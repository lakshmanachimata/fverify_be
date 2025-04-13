package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		claims, err := ParseAuthToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Check if the user's role is allowed
		for _, role := range requiredRoles {
			if claims.Role == role {
				// Add user details to the context
				c.Set("user", claims)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		c.Abort()
	}
}

const apiKey = "059d987b-5f42-44df-a13c-f7042fca0bb1"       // Replace with a secure API key
const orgAPIKey = "059d987b-5f42-44df-a13c-f7042fca0bb2"    // Replace with a secure API key
const getOrgAPIKey = "059d987b-5f42-44df-a13c-f7042fca0bb3" // Replace with a secure API key

func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the API key from the header
		providedKey := c.GetHeader("X-API-Key")
		if providedKey != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func OrgAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the API key from the header
		providedKey := c.GetHeader("X-API-Key")
		if providedKey != orgAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetOrgAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the API key from the header
		providedKey := c.GetHeader("X-API-Key")
		if providedKey != getOrgAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}
