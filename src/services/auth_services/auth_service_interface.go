package auth_services

import "github.com/gin-gonic/gin"

type AuthServiceInterface interface {
	Register(email string, username string, password string) (interface{}, error)
	Login(email, password string) (gin.H, error)
	GenerateTokens(userID int64) (*TokenResult, error)
	RefreshToken(token string) (*TokenResult, error)
	// RevokeToken(token string) error
	VerifyToken(token string) (int64, error)
}
