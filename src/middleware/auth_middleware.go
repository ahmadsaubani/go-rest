package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Format: Bearer <token>
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Load JWT secret
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not set"})
			c.Abort()
			return
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not parse token claims"})
			c.Abort()
			return
		}

		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing expiration"})
			c.Abort()
			return
		}

		// Simpan user_id ke context untuk digunakan di controller
		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(userID))
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
