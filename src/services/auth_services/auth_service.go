package auth_services

import (
	"fmt"
	"gin/src/repositories/auth_repositories"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JWTSecret = func() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("Warning: JWT_SECRET environment variable not set")
	}
	return secret
}()

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
func (s *AuthService) Register(email string, username string, password string) (interface{}, error) {
	// Call repository to perform the actual registration
	response, err := s.authRepo.Register(email, username, password)
	if err != nil {
		return nil, fmt.Errorf("could not register user: %w", err)
	}

	// Return the response from repository (already processed in the repository)
	return response, nil
}

func (s *AuthService) Login(email, password string) (gin.H, error) {
	user, err := s.authRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	tokens, err := s.GenerateTokens(user.ID)
	if err != nil {
		return nil, fmt.Errorf("could not generate tokens: %w", err)
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

	if JWTSecret == "" {
		return nil, fmt.Errorf("JWT secret is not set")
	}

	accessTokenLifetime := time.Now().Add(50 * time.Minute)
	refreshTokenLifetime := time.Now().Add(24 * 24 * time.Minute)

	accessTokenString, err := s.createJWTToken(userID, accessTokenLifetime)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := s.createJWTToken(userID, refreshTokenLifetime)
	if err != nil {
		return nil, err
	}

	// Simpan ke database via repository
	err = s.authRepo.SaveTokens(userID, accessTokenString, accessTokenLifetime, refreshTokenString, refreshTokenLifetime)
	if err != nil {
		return nil, err
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
	return token.SignedString([]byte(JWTSecret))
}

func (s *AuthService) RefreshToken(refreshTokenString string) (*TokenResult, error) {
	userID, err := s.VerifyToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	refreshTokenRecord, err := s.authRepo.FindRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found")
	}
	if refreshTokenRecord.Claimed {
		return nil, fmt.Errorf("refresh token already used")
	}

	tokenResult, err := s.GenerateTokens(userID)
	if err != nil {
		return nil, err
	}

	_ = s.authRepo.MarkRefreshTokenAsUsed(refreshTokenRecord.ID)

	return tokenResult, nil
}

func (s *AuthService) VerifyToken(tokenString string) (int64, error) {
	if JWTSecret == "" {
		return 0, fmt.Errorf("JWT secret is not set")
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil || !parsedToken.Valid {
		return 0, fmt.Errorf("invalid or expired token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id missing in claims")
	}

	return int64(userIDFloat), nil
}
