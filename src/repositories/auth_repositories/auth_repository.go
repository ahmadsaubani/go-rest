package auth_repositories

import (
	"context"
	"fmt"
	"gin/src/entities/auth"
	"gin/src/entities/users"
	"gin/src/helpers"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authRepository struct{}

func NewAuthRepository() *authRepository {
	return &authRepository{}
}

// Register handles the actual logic of saving a new user to the database
func (r *authRepository) Register(ctx context.Context, email string, username string, password string) (map[string]interface{}, error) {

	// Check if the email is already in use
	if _, err := r.FindByEmail(email); err == nil {
		return nil, fmt.Errorf("email already in use %w", err)
	}

	if _, err := r.FindByUsername(username); err == nil {
		return nil, fmt.Errorf("username already in use %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	// Create the user record
	newUser := users.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}

	// Save the user in the database
	if err := helpers.InsertModel(&newUser); err != nil {
		return nil, fmt.Errorf("could not insert user: %w", err)
	}

	// Return a map with user data
	response := map[string]interface{}{
		"id":       newUser.ID,
		"email":    newUser.Email,
		"username": newUser.Username,
	}
	return response, nil
}

// FindByEmail mencari user berdasarkan email menggunakan helper
func (r *authRepository) FindByEmail(email string) (*users.User, error) {
	var user users.User
	// Menggunakan helper untuk mencari user berdasarkan email
	err := helpers.FindOneByField(&user, "email", email)
	if err != nil {
		return nil, fmt.Errorf("email not found: %w", err)

	}
	return &user, nil
}

func (r *authRepository) FindByUsername(username string) (*users.User, error) {
	var user users.User
	err := helpers.FindOneByField(&user, "username", username)
	if err != nil {
		return nil, fmt.Errorf("username not found: %w", err)
	}
	return &user, nil
}

func (r *authRepository) CreateUser(user *users.User) error {
	return helpers.InsertModel(user)
}

func (r *authRepository) SaveTokens(userID int64, accessToken string, accessExp time.Time, refreshToken string, refreshExp time.Time) error {
	access := auth.AccessToken{
		UserID:    userID,
		Token:     accessToken,
		ExpiresAt: accessExp,
	}
	if err := helpers.InsertModel(&access); err != nil {
		return fmt.Errorf("failed insert access token: %w", err)
	}

	refresh := auth.RefreshToken{
		UserID:        userID,
		AccessTokenID: access.ID,
		Token:         refreshToken,
		ExpiresAt:     refreshExp,
	}
	return helpers.InsertModel(&refresh)
}

func (r *authRepository) FindRefreshToken(token string) (*auth.RefreshToken, error) {
	var refresh auth.RefreshToken
	if err := helpers.FindOneByField(&refresh, "token", token); err != nil {
		return nil, fmt.Errorf("token not found: %w", err)
	}
	return &refresh, nil
}

func (r *authRepository) MarkRefreshTokenAsUsed(id int64) error {
	refresh := auth.RefreshToken{
		Claimed: true,
	}
	return helpers.UpdateModelByID(&refresh, id)
}

// MarkTokenAsRevoked menandai token sebagai revoked di database
func (r *authRepository) MarkTokenAsRevoked(tokenID int64) error {
	// Buat map dengan field yang ingin diupdate
	updatedFields := map[string]interface{}{
		"revoked": true, // Hanya field revoked yang diupdate
	}

	// Panggil helper untuk update berdasarkan ID dan field yang ingin diupdate
	// Kita memastikan tipe model yang digunakan eksplisit
	return helpers.UpdateModelByIDWithMap[auth.AccessToken](updatedFields, tokenID)
}

// FindTokenByUserIDAndToken mencari token berdasarkan user_id dan token string
func (r *authRepository) FindTokenByUserIDAndToken(userID int64, tokenString string) (*auth.AccessToken, error) {
	var token auth.AccessToken
	tokenString = strings.TrimSpace(tokenString)
	// Menggunakan helper untuk mencari token berdasarkan user_id dan token string
	err := helpers.FindOneByField(&token, "user_id", userID, "token", tokenString, "revoked", false)
	if err != nil {
		return nil, fmt.Errorf("token not found or already revoked: %w", err)
	}
	return &token, nil
}
