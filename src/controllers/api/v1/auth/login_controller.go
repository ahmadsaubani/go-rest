package auth

import (
	"fmt"
	"gin/src/configs/database"
	"gin/src/entities/users"
	"gin/src/helpers"
	"gin/src/services/auth_services"

	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Login handles user login and token generation
func Login(conn *database.DBConnection, authService *auth_services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input LoginRequest

		// Bind input data from request body
		if err := ctx.ShouldBind(&input); err != nil {
			helpers.ErrorResponse(fmt.Errorf("Invalid input: %w", err), ctx, http.StatusBadRequest)

			return
		}

		// Find user by email
		var user users.User
		if err := helpers.FindOneByField(&user, "email", input.Email); err != nil {
			helpers.ErrorResponse(fmt.Errorf("Invalid email: %w", err), ctx, http.StatusUnauthorized)

			return
		}

		// Compare password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			helpers.ErrorResponse(fmt.Errorf("Invalid password"), ctx, http.StatusUnauthorized)
			return
		}

		// Generate Access and Refresh Tokens
		tokenResult, err := authService.GenerateTokens(user.ID)
		if err != nil {
			helpers.ErrorResponse(fmt.Errorf("Could not generate tokens: %w", err), ctx, http.StatusInternalServerError)
			return
		}

		// Respond with tokens
		response := gin.H{
			"token_type":         "Bearer",
			"access_token":       tokenResult.AccessToken,
			"access_expires_at":  tokenResult.AccessExpiresAt,
			"refresh_token":      tokenResult.RefreshToken,
			"refresh_expires_at": tokenResult.RefreshExpiresAt,
		}

		helpers.SuccessResponse(ctx, "Data Found!", response)
	}
}

// RefreshToken handles refreshing of access and refresh tokensc
func RefreshToken(conn *database.DBConnection, authService *auth_services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body struct {
			RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
		}
		if err := ctx.ShouldBind(&body); err != nil {
			helpers.ErrorResponse(fmt.Errorf("Invalid input: %w", err), ctx, http.StatusBadRequest)
			return
		}

		// Call the RefreshToken method in the AuthService
		tokenResult, err := authService.RefreshToken(body.RefreshToken)
		if err != nil {
			helpers.ErrorResponse(fmt.Errorf("Could not generate tokens: %w", err), ctx, http.StatusInternalServerError)
			return
		}
		helpers.SuccessResponse(ctx, "Token refreshed", tokenResult)
	}
}
