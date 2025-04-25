package auth_services

import (
	"context"

	"github.com/gin-gonic/gin"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, email string, username string, password string) (map[string]interface{}, error)
	Login(ctx context.Context, email string, password string) (gin.H, error)
	GenerateTokens(userID int64) (*TokenResult, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*TokenResult, error)
	VerifyToken(token string) (int64, error)
	RevokeToken(ctx context.Context, tokenString string) error
}
