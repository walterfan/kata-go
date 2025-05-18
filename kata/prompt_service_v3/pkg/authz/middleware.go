// pkg/authz/middleware.go
package authz

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User not authenticated"})
			return
		}

		sub := user.(string)           // e.g., "alice"
		obj := c.Request.URL.Path      // e.g., "/api/v1/prompts/1"
		act := c.Request.Method        // e.g., "GET"

		if allowed, _ := Enforcer.Enforce(sub, obj, act); allowed {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		}
	}
}