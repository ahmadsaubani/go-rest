package services

import (
	"gin/src/entities/users"
	repositories "gin/src/repositories/user_repositories"
)

type UserService interface {
	GetPaginatedUsers(limit, offset int, order string) ([]users.User, int64, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetPaginatedUsers(limit, offset int, order string) ([]users.User, int64, error) {
	usersList, err := s.repo.GetAll(limit, offset, order)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountAll()
	if err != nil {
		return nil, 0, err
	}

	return usersList, total, nil
}
