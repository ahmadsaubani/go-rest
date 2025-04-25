package services

import (
	"context"
	"fmt"
	"gin/src/entities/users"
	repositories "gin/src/repositories/user_repositories"
	"gin/src/utils/uploaders"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetPaginatedUsers(ctx *gin.Context, limit int, offset int) ([]users.User, int64, error)
	UpdateAvatar(ctx context.Context, userID int64, avatarURL string) error
	UploadAvatar(ctx *gin.Context, userID int64, file multipart.File, folder string) (string, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetPaginatedUsers(ctx *gin.Context, limit int, offset int) ([]users.User, int64, error) {
	usersList, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("sorry, we encountered an issue fetching the user list. Please try again later: %w", err)
	}

	total, err := s.repo.CountAll()
	if err != nil {
		return nil, 0, fmt.Errorf("sorry, we couldn't count the users at the moment. Please try again later: %w", err)
	}

	return usersList, total, nil
}

func (s *userService) UpdateAvatar(ctx context.Context, userID int64, avatarURL string) error {
	// Validasi avatar URL jika perlu
	if avatarURL == "" {
		return fmt.Errorf("avatar URL cannot be empty")
	}

	// Update avatar via repository
	if err := s.repo.UpdateAvatar(ctx, userID, avatarURL); err != nil {
		return fmt.Errorf("failed to update avatar: %w", err)
	}

	return nil
}

func (service *userService) UploadAvatar(ctx *gin.Context, userID int64, file multipart.File, folder string) (string, error) {
	avatarURL, err := uploaders.UploadFile(ctx.Request, folder)

	if err != nil {
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}

	if err := service.UpdateAvatar(ctx, userID, avatarURL); err != nil {
		return "", fmt.Errorf("failed to update avatar URL in database: %w", err)
	}

	// Kembalikan URL avatar yang telah berhasil disimpan
	return avatarURL, nil
}
