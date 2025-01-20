package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"poultry-management.com/internal/auth"
)

func authMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("AuthHandlerorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "AuthHandlerorization header is missing"})
			return
		}

		claims, err := auth.ValidateJWT(tokenString, jwtKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func authLocationMiddelware(jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("AuthHandlerorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "AuthHandlerorization header is missing"})
			return
		}

		claims, err := auth.ValidateJWT(tokenString, jwtKey)
		if err != nil || !(claims.Role == "admin" || claims.Role == "master" || claims.Role == "super_admin") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token/user"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}
