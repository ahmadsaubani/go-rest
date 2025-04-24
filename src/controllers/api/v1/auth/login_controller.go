package auth

import (
	"fmt"
	"gin/src/helpers"
	"gin/src/services/auth_services"

	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}

// Login memanggil fungsi ke service
func Login(authService auth_services.AuthServiceInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input LoginRequest
		if err := ctx.ShouldBind(&input); err != nil {
			helpers.ErrorResponse(fmt.Errorf("Invalid input: %w", err), ctx, http.StatusBadRequest)
			return
		}

		// Memanggil service untuk login
		response, err := authService.Login(input.Email, input.Password)
		if err != nil {
			helpers.ErrorResponse(err, ctx, http.StatusUnauthorized)
			return
		}

		helpers.SuccessResponse(ctx, "Login successful", response)
	}
}

// RefreshToken memanggil fungsi untuk menghasilkan token baru berdasarkan refresh token
func RefreshToken(authService auth_services.AuthServiceInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body RefreshTokenRequest
		if err := ctx.ShouldBind(&body); err != nil {
			helpers.ErrorResponse(fmt.Errorf("Invalid input: %w", err), ctx, http.StatusBadRequest)
			return
		}

		// Memanggil service untuk refresh token
		tokenResult, err := authService.RefreshToken(body.RefreshToken)
		if err != nil {
			helpers.ErrorResponse(fmt.Errorf("%w", err), ctx, http.StatusInternalServerError)
			return
		}

		helpers.SuccessResponse(ctx, "Token refreshed", tokenResult)
	}
}
