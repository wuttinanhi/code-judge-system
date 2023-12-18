package repositories

import (
	"github.com/wuttinanhi/code-judge-system/entities"
	repositories "github.com/wuttinanhi/code-judge-system/repositories/interfaces"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// CreateUser implements repositories.UserRepository.
func (r *userRepository) CreateUser(user *entities.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteUser implements repositories.UserRepository.
func (r *userRepository) DeleteUser(user *entities.User) error {
	result := r.db.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindUserByEmail implements repositories.UserRepository.
func (r *userRepository) FindUserByEmail(email string) (user *entities.User, err error) {
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// FindUserByID implements repositories.UserRepository.
func (r *userRepository) FindUserByID(id uint) (user *entities.User, err error) {
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// FindUserByDisplayName implements repositories.UserRepository.
func (r *userRepository) FindUserByDisplayName(displayName string) (user *entities.User, err error) {
	result := r.db.Where("display_name = ?", displayName).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// UpdateUser implements repositories.UserRepository.
func (r *userRepository) UpdateUser(user *entities.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func NewSQLUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{
		db: db,
	}
}
