package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// repository stores and controls Users in the database.
type userRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewUserRepository creates and returns a new UserRepository.
func NewUserRepository(db *gorm.DB, logger *zerolog.Logger) models.UserRepository {
	return &userRepository{db, logger}
}

// FindByID finds a single User by ID.
func (r *userRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

// FindByEmail finds a single User by Email.
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

// Create creates a new User.
func (r *userRepository) Create(user models.User) error {
	return r.db.Create(&user).Error
}

// Update updates a User with the ID in the provided User.
func (r *userRepository) Update(user models.User) error {
	return r.db.Model(&models.User{}).Updates(user).Error
}

// DeleteByID deletes a User by ID.
func (r *userRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.User{
		Model: models.Model{
			ID: id,
		},
	}).Error
}
