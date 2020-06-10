package users

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// UserRepository represents a repository of users.
type UserRepository interface {
	// FindByID finds a single User by ID.
	FindByID(id int64) (*models.User, error)
	// Create creates a new User.
	Create(user models.User) error
	// Update updates a User with the ID in the provided User.
	Update(user models.User) error
	// DeleteByID deletes a User by ID.
	DeleteByID(id int64) error
}

// repository stores and controls Users in the database.
type repository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewUserRepository creates and returns a new UserRepository.
func NewUserRepository(db *gorm.DB, logger *zerolog.Logger) UserRepository {
	return &repository{db, logger}
}

// FindByID finds a single User by ID.
func (r *repository) FindByID(id int64) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

// Create creates a new User.
func (r *repository) Create(user models.User) error {
	return r.db.Create(&user).Error
}

// Update updates a User with the ID in the provided User.
func (r *repository) Update(user models.User) error {
	return r.db.Save(&user).Error
}

// DeleteByID deletes a User by ID.
func (r *repository) DeleteByID(id int64) error {
	return r.db.Delete(&models.User{
		ID: id,
	}).Error
}
