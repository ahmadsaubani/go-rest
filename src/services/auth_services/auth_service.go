package auth_services

import (
	"fmt"
	"gin/src/entities/auth"
	"gin/src/helpers"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles auth-related operations
type AuthService struct{}

// NewAuthService initializes a new AuthService
func NewAuthService() *AuthService {
	return &AuthService{}
}

type TokenResult struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

// GenerateTokens generates access and refresh tokens and stores them in the database
func (s *AuthService) GenerateTokens(userID int64) (*TokenResult, error) {
	fmt.Println("Generating tokens for user ID:", userID)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT secret is not set")
	}

	accessTokenLifetime := time.Now().Add(50 * time.Minute)
	refreshTokenLifetime := time.Now().Add(24 * 24 * time.Minute) // 24 hari

	// Access Token
	accessTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     accessTokenLifetime.Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	accessTokenRecord := auth.AccessToken{
		UserID:    userID,
		Token:     accessTokenString,
		ExpiresAt: accessTokenLifetime,
	}
	if err := helpers.InsertModel(&accessTokenRecord); err != nil {
		return nil, err
	}

	// Refresh Token
	refreshTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     refreshTokenLifetime.Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	refreshTokenRecord := auth.RefreshToken{
		UserID:        userID,
		AccessTokenID: accessTokenRecord.ID,
		Token:         refreshTokenString,
		ExpiresAt:     refreshTokenLifetime,
	}
	if err := helpers.InsertModel(&refreshTokenRecord); err != nil {
		return nil, err
	}

	return &TokenResult{
		AccessToken:      accessTokenString,
		RefreshToken:     refreshTokenString,
		AccessExpiresAt:  accessTokenLifetime,
		RefreshExpiresAt: refreshTokenLifetime,
	}, nil
}

// RefreshToken verifies and regenerates tokens if the refresh token is valid
func (s *AuthService) RefreshToken(refreshTokenString string) (*TokenResult, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT secret is not set")
	}

	parsedToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims in token")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("user_id is missing or invalid in token claims")
	}
	userID := int64(userIDFloat)

	var refreshTokenRecord auth.RefreshToken
	err = helpers.FindOneByField(&refreshTokenRecord, "token", refreshTokenString)
	fmt.Println("Refresh token record:", err)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found")
	}

	if refreshTokenRecord.Claimed {
		return nil, fmt.Errorf("refresh token already used")
	}

	// Generate new tokens
	tokenResult, err := s.GenerateTokens(userID)
	if err != nil {
		return nil, err
	}

	// Tandai refresh token lama sebagai digunakan
	refreshTokenRecord.Claimed = true
	if err := helpers.UpdateModelByID(&refreshTokenRecord, refreshTokenRecord.ID); err != nil {
		return nil, err
	}

	return tokenResult, nil
}

// RevokeToken marks a refresh token as claimed (used)
func (s *AuthService) RevokeToken(token string) error {
	var refreshToken auth.RefreshToken
	err := helpers.FindOneByField(&refreshToken, "token", token)
	if err != nil {
		return fmt.Errorf("token not found: %v", err)
	}

	refreshToken.Claimed = true
	return helpers.UpdateModelByID(&refreshToken, refreshToken.ID)
}

// VerifyToken verifies a JWT access token and returns the user ID if valid
func (s *AuthService) VerifyToken(tokenString string) (uint, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return 0, fmt.Errorf("JWT secret is not set")
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !parsedToken.Valid {
		return 0, fmt.Errorf("invalid or expired token: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("failed to parse claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id claim is missing or invalid")
	}

	return uint(userIDFloat), nil
}
