package repositories

import (
	"gin/src/entities/users"
	"gin/src/helpers"

	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	GetAll(ctx *gin.Context, limit int, offset int, order string) ([]users.User, error)
	CountAll() (int64, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) GetAll(ctx *gin.Context, limit, offset int, order string) ([]users.User, error) {
	var usersList []users.User

	// err := helpers.GetAllModels(ctx, &usersList, limit, offset, order)
	err := helpers.GetAllModels(ctx, &usersList, limit, offset, order)

	return usersList, err
}

func (r *userRepository) CountAll() (int64, error) {
	return helpers.CountModel[users.User]()
}
