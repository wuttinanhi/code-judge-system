package repositories

import (
	"strings"

	"github.com/wuttinanhi/code-judge-system/entities"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	// FindUserByID returns a user by given ID.
	FindUserByID(id uint) (user *entities.User, err error)
	// FindUserByDisplayName returns a user by given display name.
	FindUserByDisplayName(displayName string) (user *entities.User, err error)
	// FindUserByEmail returns a user by given email.
	FindUserByEmail(email string) (user *entities.User, err error)
	// CreateUser creates a new user.
	CreateUser(user *entities.User) (*entities.User, error)
	// UpdateUser updates a user.
	UpdateUser(user *entities.User) error
	// DeleteUser deletes a user.
	DeleteUser(user *entities.User) error
	// Pagination returns a list of users by given page and limit.
	Pagination(options *entities.PaginationOptions) (result *entities.PaginationResult[*entities.User], err error)
}

type userRepository struct {
	db *gorm.DB
}

// Pagination implements UserRepository.
func (r *userRepository) Pagination(options *entities.PaginationOptions) (result *entities.PaginationResult[*entities.User], err error) {
	offset := (options.Page - 1) * options.Limit
	desc := strings.ToUpper(options.Order) == "DESC"

	findQuery := r.db.Model(&entities.User{}).
		Where(
			"display_name LIKE ? OR email LIKE ?",
			"%"+options.Search+"%",
			"%"+options.Search+"%",
		).
		Limit(options.Limit).
		Offset(offset).
		Order(clause.OrderByColumn{Column: clause.Column{Name: options.Sort}, Desc: desc})

	var users []*entities.User
	if err := findQuery.Find(&users).Error; err != nil {
		return nil, err
	}

	var total int64
	countQuery := r.db.Model(&entities.User{}).
		Where(
			"display_name LIKE ? OR email LIKE ?",
			"%"+options.Search+"%",
			"%"+options.Search+"%",
		).Count(&total)
	if err := countQuery.Error; err != nil {
		return nil, err
	}

	return &entities.PaginationResult[*entities.User]{
		Total: int(total),
		Items: users,
	}, nil
}

// CreateUser implements repositories.UserRepository.
func (r *userRepository) CreateUser(user *entities.User) (*entities.User, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// DeleteUser implements repositories.UserRepository.
func (r *userRepository) DeleteUser(user *entities.User) error {
	result := r.db.Delete(user)
	return result.Error
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
	return result.Error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
