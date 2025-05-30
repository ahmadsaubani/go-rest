package auth_services

import (
	"context"
	"fmt"
	"gin/src/repositories/auth_repositories"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, email string, username string, password string) (map[string]interface{}, error)
	Login(ctx context.Context, email string, password string) (gin.H, error)
	GenerateTokens(userID int64) (*TokenResult, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*TokenResult, error)
	VerifyToken(token string) (int64, error)
	RevokeToken(ctx context.Context, tokenString string) error
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is not set")
	}
	return secret
}

type AuthService struct {
	authRepo auth_repositories.AuthRepositoryInterface
}

func NewAuthService(repo auth_repositories.AuthRepositoryInterface) *AuthService {
	return &AuthService{authRepo: repo}
}

type TokenResult struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

// Register handles the registration logic
func (s *AuthService) Register(ctx context.Context, email string, username string, password string) (map[string]interface{}, error) {
	response, err := s.authRepo.Register(ctx, email, username, password)
	if err != nil {
		return nil, fmt.Errorf("could not register user: %w", err)
	}
	return response, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (gin.H, error) {

	user, err := s.authRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	tokens, err := s.GenerateTokens(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed generate tokens: %w", err)
	}

	return gin.H{
		"token_type":         "Bearer",
		"access_token":       tokens.AccessToken,
		"access_expires_at":  tokens.AccessExpiresAt,
		"refresh_token":      tokens.RefreshToken,
		"refresh_expires_at": tokens.RefreshExpiresAt,
	}, nil
}

func (s *AuthService) GenerateTokens(userID int64) (*TokenResult, error) {

	accessTokenLifetime := time.Now().Add(50 * time.Minute)
	refreshTokenLifetime := time.Now().Add(24 * 24 * time.Minute)

	accessTokenString, err := s.createJWTToken(userID, accessTokenLifetime)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT token for access token: %w", err)
	}

	refreshTokenString, err := s.createJWTToken(userID, refreshTokenLifetime)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT token for refresh token: %w", err)
	}

	// Simpan ke database via repository
	err = s.authRepo.SaveTokens(userID, accessTokenString, accessTokenLifetime, refreshTokenString, refreshTokenLifetime)
	if err != nil {
		return nil, fmt.Errorf("save token to database error: %w", err)
	}

	return &TokenResult{
		AccessToken:      accessTokenString,
		RefreshToken:     refreshTokenString,
		AccessExpiresAt:  accessTokenLifetime,
		RefreshExpiresAt: refreshTokenLifetime,
	}, nil
}

func (s *AuthService) createJWTToken(userID int64, exp time.Time) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getJWTSecret()))
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (*TokenResult, error) {

	userID, err := s.VerifyToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token: %w", err)
	}

	refreshTokenRecord, err := s.authRepo.FindRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found: %w", err)

	}
	if refreshTokenRecord.Claimed {
		return nil, fmt.Errorf("refresh token already claimed and used: %w", err)
	}

	tokenResult, err := s.GenerateTokens(userID)
	if err != nil {
		return nil, fmt.Errorf("error generate tokens: %w", err)
	}

	_ = s.authRepo.MarkRefreshTokenAsUsed(refreshTokenRecord.ID)

	return tokenResult, nil
}

func (s *AuthService) VerifyToken(tokenString string) (int64, error) {

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(getJWTSecret()), nil
	})
	if err != nil || !parsedToken.Valid {
		return 0, fmt.Errorf("invalid or expired token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims format, expected jwt.MapClaims, got: %T. Claims: %v", parsedToken.Claims, claims)
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id missing in claims or not a float64, claims: %v", claims)
	}

	return int64(userIDFloat), nil
}

func (s *AuthService) RevokeToken(ctx context.Context, tokenString string) error {

	// Verifikasi apakah token yang diberikan valid
	userID, err := s.VerifyToken(tokenString)
	if err != nil {
		return fmt.Errorf("invalid or expired token: %w", err)
	}

	// Cari token dalam database
	tokenRecord, err := s.authRepo.FindTokenByUserIDAndToken(userID, tokenString)

	if err != nil {
		return fmt.Errorf("token not found: %w", err)
	}

	// Tandai token sebagai revoked
	err = s.authRepo.MarkTokenAsRevoked(tokenRecord.ID)
	if err != nil {
		return fmt.Errorf("failed to mark token as revoked: %w", err)
	}

	return nil
}
