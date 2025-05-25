package auth

import (
	"fverify_be/internal/repositories"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func AuthMiddleware(orgRepo repositories.OrganisationRepositoryImpl, userRepo repositories.UserRepositoryImpl, requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required with format Bearer <token>"})
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

		// Extract org_id from the request (assuming it's passed as a query parameter or path parameter)
		org_id := c.GetHeader("org_id")
		if org_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "org_id is required"})
			c.Abort()
			return
		}

		// Step 1: Get org from org_id
		org, err := orgRepo.GetOrganisationByID(c.Request.Context(), org_id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organisation not found"})
			c.Abort()
			return
		}

		// Step 2: Check if orgUUID from token matches orgUUID fetched from org
		if claims.OrgUUID != org.OrgUUID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied: Organisation mismatch"})
			c.Abort()
			return
		}

		// Step 3: Check if org status and user status are active
		if string(org.Status) != "Active" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Inactive organisation"})
			c.Abort()
			return
		}

		// Step 4: Get user from claims.UserId
		user, err := userRepo.GetByUserID(c.Request.Context(), claims.UserId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Step 5: Check if user status is active
		if string(user.Status) != "Active" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied: Inactive user"})
			c.Abort()
			return
		}

		// Check if the user's role is allowed
		for _, role := range requiredRoles {
			if claims.Role == role {
				// Add user and org details to the context
				c.Set("user", claims)
				c.Set("org", org)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions to access this resource"})
		c.Abort()
	}
}

func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := viper.GetString("apikeys.userAPIKey")
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
		orgAPIKey := viper.GetString("apikeys.orgAPIKey")
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
