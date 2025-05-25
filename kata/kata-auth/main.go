package main

import (
	"net/http"

	middleware "github.com/walterfan/kata-auth/internal"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	r := gin.Default()
	enforcer, _ := casbin.NewEnforcer("config/model.conf", "config/policy.csv")

	// Public endpoint to get a test token
	r.POST("/token", func(c *gin.Context) {
		var loginReq LoginRequest
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Example validation (replace with real authentication logic)
		if loginReq.Username != "test" || loginReq.Password != "pass" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "123",
			"role":    "admin", // You could also dynamically assign role
		})
		tokenStr, _ := token.SignedString([]byte("my-secret"))
		c.JSON(http.StatusOK, gin.H{"token": tokenStr})
	})

	auth := r.Group("/")
	auth.Use(middleware.JWTAuth())

	auth.GET("/admin", func(c *gin.Context) {
		role := middleware.GetUserRole(c)
		ok, _ := enforcer.Enforce(role, "/admin", "GET")
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		c.String(http.StatusOK, "Welcome Admin!")
	})

	auth.GET("/user", func(c *gin.Context) {
		role := middleware.GetUserRole(c)
		ok, _ := enforcer.Enforce(role, "/user", "GET")
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		userID := middleware.GetUserID(c)
		c.String(http.StatusOK, "Hello user %s", userID)
	})

	r.Run(":8080")
}
