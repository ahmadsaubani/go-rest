package services

import (
	"fmt"
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
		return nil, 0, fmt.Errorf("sorry, we encountered an issue fetching the user list. Please try again later: %w", err)
	}

	total, err := s.repo.CountAll()
	if err != nil {
		return nil, 0, fmt.Errorf("sorry, we couldn't count the users at the moment. Please try again later: %w", err)
	}

	return usersList, total, nil
}
