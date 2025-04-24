package auth_repositories

import (
	"context"
	"gin/src/entities/auth"
	"gin/src/entities/users"
	"time"
)

type AuthRepositoryInterface interface {
	Register(ctx context.Context) (map[string]interface{}, error)
	FindByEmail(email string) (*users.User, error)
	FindByUsername(username string) (*users.User, error)
	CreateUser(user *users.User) error
	SaveTokens(userID int64, accessToken string, accessExp time.Time, refreshToken string, refreshExp time.Time) error
	FindRefreshToken(token string) (*auth.RefreshToken, error)
	MarkRefreshTokenAsUsed(id int64) error
}
