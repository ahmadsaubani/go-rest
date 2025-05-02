package registrations

import (
	"gin/src/middleware"

	"github.com/gin-gonic/gin"
)

func GlobalMiddlewares(ginEngine *gin.Engine) *gin.Engine {
	// untuk rate limiter
	ginEngine.Use(middleware.RateLimiter())
	// untuk cors
	ginEngine.Use(middleware.SecureHeadersMiddleware())
	// untuk logger
	ginEngine.Use(middleware.RecoveryWithLogger())
	// untuk request body
	ginEngine.Use(middleware.SaveRequestBody())

	return ginEngine
}
