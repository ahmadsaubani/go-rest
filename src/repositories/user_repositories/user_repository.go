package repositories

import (
	"fmt"
	"gin/src/entities/users"
	"gin/src/helpers"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type UserRepository interface {
	GetAll(ctx *gin.Context, limit int, offset int) ([]users.User, error)
	CountAll() (int64, error)
	UpdateAvatar(ctx context.Context, userID int64, avatarURL string) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) GetAll(ctx *gin.Context, limit, offset int) ([]users.User, error) {
	var usersList []users.User
	err := helpers.GetAllModels(ctx, &usersList, limit, offset)

	return usersList, err
}

func (r *userRepository) CountAll() (int64, error) {
	return helpers.CountModel[users.User]()
}

func (r *userRepository) UpdateAvatar(ctx context.Context, userID int64, avatarURL string) error {
	var user users.User
	if err := helpers.GetModelByID(&user, userID); err != nil {
		return fmt.Errorf("failed to find user by ID: %w", err)
	}
	// Update avatar URL

	// Simpan perubahan
	updatedFields := map[string]interface{}{
		"avatar": avatarURL,
	}

	// Panggil helper untuk update berdasarkan ID dan field yang ingin diupdate
	// Kita memastikan tipe model yang digunakan eksplisit
	return helpers.UpdateModelByIDWithMap[users.User](updatedFields, userID)
}
