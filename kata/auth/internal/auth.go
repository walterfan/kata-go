package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "my-secret"

// Middleware to parse and validate JWT, and store claims in context
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, "Bearer ")
        if len(parts) != 2 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
            c.Abort()
            return
        }

        tokenStr := parts[1]
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
            c.Abort()
            return
        }

        c.Set("userID", claims["user_id"])
        c.Set("role", claims["role"])
        c.Next()
    }
}

// GetUserRole returns role from Gin context
func GetUserRole(c *gin.Context) string {
    if val, exists := c.Get("role"); exists {
        return val.(string)
    }
    return ""
}

// GetUserID returns user ID from Gin context
func GetUserID(c *gin.Context) string {
    if val, exists := c.Get("userID"); exists {
        return val.(string)
    }
    return ""
}
