package repositories

import (
	"gin/src/entities/users"
	"gin/src/helpers"
)

type UserRepository interface {
	GetAll(limit, offset int, order string) ([]users.User, error)
	CountAll() (int64, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) GetAll(limit, offset int, order string) ([]users.User, error) {
	var usersList []users.User
	err := helpers.GetAllModels(&usersList, limit, offset, order)
	return usersList, err
}

func (r *userRepository) CountAll() (int64, error) {
	return helpers.CountModel[users.User]()
}
