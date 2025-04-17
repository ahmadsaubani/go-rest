package routes

import (
	"gin/src/http/controllers/api/v1/auth"
	"gin/src/http/controllers/api/v1/user"
	"gin/src/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func API() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		v1.POST("/user/register", auth.Register)
		v1.POST("/user/login", auth.Login)

		v1.Use(middleware.JWTAuthMiddleware()) // 🔐 Apply middleware here
		{
			v1.GET("/user/profile", user.GetProfile)
		}
	}

	return r
}
