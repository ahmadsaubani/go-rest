package routes

import (
	"gin/src/configs/database"
	"gin/src/controllers/api/v1/auth"
	"gin/src/controllers/api/v1/user"
	"gin/src/middleware"
	"gin/src/services/auth_services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func API(db *database.DBConnection) *gin.Engine {

	r := gin.Default()

	// Initialize the AuthService
	authService := auth_services.NewAuthService()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		v1.POST("/user/register", auth.Register)
		v1.POST("/user/login", auth.Login(db, authService))

		v1.Use(middleware.JWTAuthMiddleware())
		{
			v1.GET("/user/profile", user.GetProfile)
			v1.GET("/users", user.GetAllUsers)

			v1.POST("/token/refresh", auth.RefreshToken(db, authService))
		}
	}

	return r
}
