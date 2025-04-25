package auth

import (
	"fmt"
	"gin/src/helpers"
	"gin/src/services/auth_services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Logout(authService auth_services.AuthServiceInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Ambil refresh token dari header Authorization (Bearer token)
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			helpers.ErrorResponse(ctx, fmt.Errorf("authorization token missing in header"), http.StatusUnauthorized)
			return
		}

		// Ekstrak token dari Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Panggil service untuk revoke tokens
		err := authService.RevokeToken(ctx.Request.Context(), tokenString)
		if err != nil {
			helpers.ErrorResponse(ctx, err, http.StatusBadRequest)
			return
		}
		helpers.SuccessResponse(ctx, "Token revoked successfully", nil)
	}
}
