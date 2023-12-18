package services

import (
	"errors"
	"strings"

	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/repositories"
)

type UserService interface {
	Register(email, password, displayname string) (user *entities.User, err error)
	Login(email, password string) (user *entities.User, err error)
}

type userService struct {
	userRepo repositories.UserRepository
}

// Login implements services.UserService.
func (s *userService) Login(email string, password string) (user *entities.User, err error) {
	user, err = s.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("failed to get user")
	}

	if user.Password != password {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

// Register implements services.UserService.
func (s *userService) Register(email string, password string, displayname string) (user *entities.User, err error) {
	user = &entities.User{
		Email:       email,
		Password:    password,
		DisplayName: displayname,
		Role:        "user",
	}

	// Save the user to the repository
	err = s.userRepo.CreateUser(user)
	if err != nil {
		// if error contains "UNIQUE constraint failed"
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, errors.New("user already exists")
		}

		return nil, errors.New("failed to register user")
	}

	return user, nil
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
