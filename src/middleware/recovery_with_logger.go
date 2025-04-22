package middleware

import (
	"fmt"
	"gin/src/utils/loggers"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

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
