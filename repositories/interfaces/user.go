package repositories

import "github.com/wuttinanhi/code-judge-system/entities"

type UserRepository interface {
	// FindUserByID returns a user by given ID.
	FindUserByID(id uint) (user *entities.User, err error)
	// FindUserByDisplayName returns a user by given display name.
	FindUserByDisplayName(displayName string) (user *entities.User, err error)
	// FindUserByEmail returns a user by given email.
	FindUserByEmail(email string) (user *entities.User, err error)
	// CreateUser creates a new user.
	CreateUser(user *entities.User) error
	// UpdateUser updates a user.
	UpdateUser(user *entities.User) error
	// DeleteUser deletes a user.
	DeleteUser(user *entities.User) error
}
