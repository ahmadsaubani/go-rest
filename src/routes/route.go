package routes

import (
	"gin/src/configs/database"
	"gin/src/controllers/api/v1/auth"
	"gin/src/controllers/api/v1/user"
	"gin/src/middleware"
	"gin/src/repositories/auth_repositories"
	repositories "gin/src/repositories/user_repositories"
	"gin/src/services/auth_services"
	services "gin/src/services/user_services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API sets up the routes and handlers for the application.
// It initializes the authentication and user services using their respective repositories.
// The function defines a versioned API group (/api/v1) and registers various endpoints:
// - GET /ping: Responds with a "pong" message for health checks.
// - POST /user/register: Registers a new user using the provided authentication service.
// - POST /user/login: Authenticates a user with the provided credentials.
// - Secures routes with JWT middleware, ensuring protected endpoints require valid tokens:
//   - GET /user/profile: Returns the profile of the authenticated user.
//   - GET /users: Retrieves a list of users using the user service.
//   - POST /user/upload/avatar: Allows users to upload avatars.
//   - POST /token/refresh: Refreshes JWT tokens.
//   - POST /user/logout: Logs out the user, revoking the current token.
// Returns the configured Gin engine instance.

func API(db *database.DBConnection, ginEngine *gin.Engine) *gin.Engine {
	authRepo := auth_repositories.NewAuthRepository()
	authService := auth_services.NewAuthService(authRepo)

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)

	v1 := ginEngine.Group("/api/v1")
	{
		v1.GET("/ping", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		v1.POST("/user/register", auth.Register(authService))
		v1.POST("/user/login", auth.Login(authService))

		v1.Use(middleware.JWTAuthMiddleware())
		{
			v1.GET("/user/profile", user.GetProfile)
			v1.GET("/users", user.GetAllUsers(userService))
			v1.POST("/user/upload/avatar", user.UploadAvatar(userService))

			v1.POST("/token/refresh", auth.RefreshToken(authService))
			v1.POST("/user/logout", auth.Logout(authService))
		}
	}

	return ginEngine
}
