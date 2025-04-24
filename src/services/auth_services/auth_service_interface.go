package auth_services

import (
	"context"

	"github.com/gin-gonic/gin"
)

type AuthServiceInterface interface {
	Register(ctx context.Context) (interface{}, error)
	Login(ctx context.Context) (gin.H, error)
	GenerateTokens(userID int64) (*TokenResult, error)
	RefreshToken(ctx context.Context) (*TokenResult, error)
	VerifyToken(token string) (int64, error)
}
