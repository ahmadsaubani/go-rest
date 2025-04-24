package auth

import (
	"gin/src/helpers"
	"gin/src/services/auth_services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
}

func Register(authService auth_services.AuthServiceInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body RegisterRequest

		// Bind the request body to struct
		if err := ctx.ShouldBind(&body); err != nil {
			helpers.ErrorResponse(err, ctx, http.StatusBadRequest)
			return
		}

		// Call service to register the user
		response, err := authService.Register(body.Email, body.Username, body.Password)
		if err != nil {
			helpers.ErrorResponse(err, ctx, http.StatusBadRequest)
			return
		}

		// Return the response from service
		helpers.SuccessResponse(ctx, "User registered successfully", response)
	}
}
