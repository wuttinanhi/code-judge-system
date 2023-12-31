package services

import (
	"errors"
	"strings"

	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(email, password, displayname string) (user *entities.User, err error)
	Login(email, password string) (user *entities.User, err error)
	UpdateRole(user *entities.User, role string) (err error)
	FindUserByID(id uint) (user *entities.User, err error)
	Pagination(options *entities.PaginationOptions) (result *entities.PaginationResult[*entities.User], err error)
}

type userService struct {
	userRepo repositories.UserRepository
}

// Pagination implements UserService.
func (s *userService) Pagination(options *entities.PaginationOptions) (result *entities.PaginationResult[*entities.User], err error) {
	result, err = s.userRepo.Pagination(options)
	return
}

// FindUserByID implements UserService.
func (s *userService) FindUserByID(id uint) (user *entities.User, err error) {
	user, err = s.userRepo.FindUserByID(id)
	return user, err
}

// UpdateRole implements UserService.
func (s *userService) UpdateRole(user *entities.User, role string) (err error) {
	user.Role = role
	err = s.userRepo.UpdateUser(user)
	return err
}

// Login implements services.UserService.
func (s *userService) Login(email string, password string) (user *entities.User, err error) {
	user, err = s.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("failed to get user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

// Register implements services.UserService.
func (s *userService) Register(email string, password string, displayname string) (user *entities.User, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user = &entities.User{
		Email:       email,
		Password:    string(hashedPassword),
		DisplayName: displayname,
		Role:        "user",
	}

	// Save the user to the repository
	user, err = s.userRepo.CreateUser(user)
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
