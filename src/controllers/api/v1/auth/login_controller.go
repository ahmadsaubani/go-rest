package auth

import (
	"gin/src/entities/users"
	"gin/src/services/auth_services"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// func Login(ctx *gin.Context) {
// 	var input struct {
// 		Email    string `json:"email" binding:"required,email"`
// 		Password string `json:"password" binding:"required"`
// 	}

// 	// Bind input
// 	if err := ctx.ShouldBindJSON(&input); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	// Find user
// 	var user users.User
// 	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
// 		return
// 	}

// 	// Compare password
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
// 		return
// 	}

// 	// JWT Secret
// 	secret := os.Getenv("JWT_SECRET")
// 	// refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
// 	// if secret == "" || refreshSecret == "" {
// 	if secret == "" {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secrets not set"})
// 		return
// 	}

// 	// Access token (24h)
// 	accessExp := time.Now().Add(24 * time.Hour)
// 	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"user_id": user.ID,
// 		"exp":     accessExp.Unix(),
// 	})

// 	accessTokenString, err := accessToken.SignedString([]byte(secret))
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
// 		return
// 	}

// 	// Refresh token (7 days)
// 	refreshExp := time.Now().Add(7 * 24 * time.Hour)
// 	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"user_id": user.ID,
// 		"exp":     refreshExp.Unix(),
// 	})

// 	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
// 		return
// 	}

// 	refreshExpires := time.Now().Add(7 * 24 * time.Hour) // Refresh token valid for 7 days

// 	// Save the refresh token to DB or in-memory store with expiry and user ID
// 	database.DB.Create(&auth.RefreshToken{
// 		UserID:    uint(user.ID),
// 		Token:     refreshTokenString,
// 		ExpiresAt: refreshExpires,
// 	})

// 	// Success response
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "Login successful",
// 		"data": gin.H{
// 			"access_token":       accessTokenString,
// 			"expires_at":         accessExp.Format(time.RFC3339),
// 			"refresh_token":      refreshTokenString,
// 			"refresh_expires_at": refreshExp.Format(time.RFC3339),
// 			"token_type":         "Bearer",
// 		},
// 	})
// }

// func RefreshToken(ctx *gin.Context) {
// 	var input struct {
// 		RefreshToken string `json:"refresh_token" binding:"required"`
// 	}

// 	if err := ctx.ShouldBindJSON(&input); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
// 		return
// 	}

// 	secret := os.Getenv("JWT_SECRET")
// 	if secret == "" {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not set"})
// 		return
// 	}

// 	// Parse and validate the token
// 	parsedToken, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method")
// 		}
// 		return []byte(secret), nil
// 	})

// 	if err != nil || !parsedToken.Valid {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
// 		return
// 	}

// 	claims, ok := parsedToken.Claims.(jwt.MapClaims)
// 	if !ok || !parsedToken.Valid {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
// 		return
// 	}

// 	userIDFloat, ok := claims["user_id"].(float64)
// 	if !ok {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
// 		return
// 	}
// 	userID := int64(userIDFloat)

// 	// Check if the refresh token exists, not expired, and not used (claimed)
// 	var storedToken auth.RefreshToken
// 	err = database.DB.
// 		Where("token = ? AND user_id = ? AND claimed = ?", input.RefreshToken, userID, false).
// 		First(&storedToken).Error

// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found or already used"})
// 		return
// 	}

// 	// Mark the old refresh token as claimed
// 	storedToken.Claimed = true
// 	database.DB.Save(&storedToken)

// 	// Create new access token
// 	expiration := time.Now().Add(24 * time.Hour)
// 	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"user_id": userID,
// 		"exp":     expiration.Unix(),
// 	})
// 	accessTokenString, err := accessToken.SignedString([]byte(secret))
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign access token"})
// 		return
// 	}

// 	// Create new refresh token
// 	newRefreshExp := time.Now().Add(7 * 24 * time.Hour)
// 	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"user_id": userID,
// 		"exp":     newRefreshExp.Unix(),
// 	})
// 	newRefreshTokenString, err := newRefreshToken.SignedString([]byte(secret))
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign refresh token"})
// 		return
// 	}

// 	// Store new refresh token
// 	newTokenRecord := auth.RefreshToken{
// 		UserID:    uint(userID),
// 		Token:     newRefreshTokenString,
// 		ExpiresAt: newRefreshExp,
// 		Claimed:   false,
// 	}
// 	database.DB.Create(&newTokenRecord)

// 	// Return new tokens
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "Token refreshed",
// 		"data": gin.H{
// 			"access_token":  accessTokenString,
// 			"refresh_token": newRefreshTokenString,
// 			"expires_in":    expiration.Unix(),
// 			"token_type":    "Bearer",
// 		},
// 	})
// }

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login handles user login and token generation
func Login(db *gorm.DB, authService *auth_services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input LoginRequest

		// Bind input data from request body
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Find user by email
		var user users.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Compare password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Generate Access and Refresh Tokens
		accessToken, refreshToken, err := authService.GenerateTokens(uint(user.ID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate tokens"})
			return
		}

		// Respond with tokens
		ctx.JSON(http.StatusOK, gin.H{
			"success":       true,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    time.Now().Add(24 * time.Hour).Unix(),
		})
	}
}

// RefreshToken handles refreshing of access and refresh tokens
func RefreshToken(db *gorm.DB, authService *auth_services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Call the RefreshToken method in the AuthService
		accessToken, newRefreshToken, err := authService.RefreshToken(body.RefreshToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}

		// Respond with new tokens
		ctx.JSON(http.StatusOK, gin.H{
			"success":       true,
			"access_token":  accessToken,
			"refresh_token": newRefreshToken,
			"token_type":    "Bearer",
			"expires_in":    time.Now().Add(24 * time.Hour).Unix(),
		})
	}
}
