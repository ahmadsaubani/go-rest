package middleware

import (
	"fmt"
	"gin/src/utils/loggers"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// RecoveryWithLogger returns a middleware that recover from any panics and logs the error to the configured logger.
// The error response will be sent to the client with a 500 status code and a JSON response with "success" set to false and "message" set to "Internal server error".
func RecoveryWithLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Log error ke file
				loggers.Log.Error("Panic Recovered", map[string]interface{}{
					"error":     fmt.Sprintf("%v", r),
					"path":      c.FullPath(),
					"method":    c.Request.Method,
					"client_ip": c.ClientIP(),
					"stack":     string(debug.Stack()),
				})

				// Respond error ke client
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Internal server error",
				})
			}
		}()

		c.Next()
	}
}
