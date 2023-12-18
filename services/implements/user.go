package services

import (
	"errors"

	"github.com/wuttinanhi/code-judge-system/entities"
	repositories "github.com/wuttinanhi/code-judge-system/repositories/interfaces"
	services "github.com/wuttinanhi/code-judge-system/services/interfaces"
)

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
		return nil, errors.New("failed to register user")
	}

	return user, nil
}

func NewUserService(userRepo repositories.UserRepository) services.UserService {
	return &userService{
		userRepo: userRepo,
	}
}
