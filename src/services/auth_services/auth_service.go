package auth_services

import (
	"fmt"
	"gin/src/entities/auth"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// AuthService is a struct for the token generation and refresh logic
type AuthService struct {
	DB *gorm.DB
}

// NewAuthService initializes a new AuthService
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

// GenerateTokens generates access and refresh tokens and stores them in the database
func (s *AuthService) GenerateTokens(userID uint) (string, string, error) {
	// Get the JWT secret from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", "", fmt.Errorf("JWT secret is not set")
	}

	accessTokenLifetime := time.Now().Add(1 * time.Minute)
	// Generate Access Token
	accessTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     accessTokenLifetime.Unix(), // expired in 5 minutes

	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	// Store refresh token in the database (separate table for refresh tokens)
	accessTokenRecord := auth.AccessToken{
		UserID:    userID,
		Token:     accessTokenString,
		ExpiresAt: accessTokenLifetime,
	}

	// Save refresh token record to database
	if err := s.DB.Create(&accessTokenRecord).Error; err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days for refresh token
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	// Store refresh token in the database (separate table for refresh tokens)
	refreshTokenRecord := auth.RefreshToken{
		UserID:        userID,
		AccessTokenID: accessTokenRecord.ID,
		Token:         refreshTokenString,
		ExpiresAt:     time.Now().Add(30 * 24 * time.Hour),
	}

	// Save refresh token record to database
	if err := s.DB.Create(&refreshTokenRecord).Error; err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// RefreshToken generates a new access token and refresh token when the refresh token is valid
func (s *AuthService) RefreshToken(refreshTokenString string) (string, string, error) {
	// Load JWT secret
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", "", fmt.Errorf("JWT secret is not set")
	}

	// Debug log
	fmt.Println("🔍 Incoming refresh token:", refreshTokenString)

	// Parse token
	parsedToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", "", fmt.Errorf("token parse error: %v", err)
	}

	if !parsedToken.Valid {
		return "", "", fmt.Errorf("invalid or expired refresh token")
	}

	// Extract claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return "", "", fmt.Errorf("could not extract claims from token")
	}

	// Parse user ID safely (JWT stores numeric values as float64)
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return "", "", fmt.Errorf("user_id is missing or invalid in token claims")
	}
	userID := uint(userIDFloat)

	// Validate token exists and not claimed
	var refreshTokenRecord auth.RefreshToken
	if err := s.DB.Where("token = ? AND claimed = ?", refreshTokenString, false).First(&refreshTokenRecord).Error; err != nil {
		return "", "", fmt.Errorf("refresh token not found or already used")
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := s.GenerateTokens(userID)
	if err != nil {
		return "", "", err
	}

	// Mark current refresh token as claimed
	refreshTokenRecord.Claimed = true
	if err := s.DB.Save(&refreshTokenRecord).Error; err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}
