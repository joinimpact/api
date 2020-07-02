package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// userProfileFieldRepository is the postgres implementation of UserProfileFieldRepository.
type userProfileFieldRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewUserProfileFieldRepository creates and returns a new UserProfileFieldRepository.
func NewUserProfileFieldRepository(db *gorm.DB, logger *zerolog.Logger) models.UserProfileFieldRepository {
	return &userProfileFieldRepository{
		db,
		logger,
	}
}

// FindByID finds a single entity by ID.
func (r *userProfileFieldRepository) FindByID(id int64) (*models.UserProfileField, error) {
	var profileField models.UserProfileField
	if err := r.db.First(&profileField, id).Error; err != nil {
		return &profileField, err
	}
	return &profileField, nil
}

// FindByUserID finds entities by UserID.
func (r *userProfileFieldRepository) FindByUserID(id int64) ([]models.UserProfileField, error) {
	var profileFields []models.UserProfileField
	if err := r.db.Where("user_id = ?", id).Find(&profileFields).Error; err != nil {
		return profileFields, err
	}
	return profileFields, nil
}

// FindUserFieldByName finds a single entity by UserID and field name.
func (r *userProfileFieldRepository) FindUserFieldByName(id int64, name string) (*models.UserProfileField, error) {
	var profileField models.UserProfileField
	if err := r.db.Where("user_id = ? AND name = ?", id, name).Find(&profileField).Error; err != nil {
		return &profileField, err
	}
	return &profileField, nil
}

// Create creates a new entity.
func (r *userProfileFieldRepository) Create(profileField models.UserProfileField) error {
	return r.db.Create(&profileField).Error
}

// Update updates an entity with the ID in the provided entity.
func (r *userProfileFieldRepository) Update(profileField models.UserProfileField) error {
	return r.db.Model(&models.UserProfileField{}).Updates(&profileField).Error
}

// DeleteByID deletes an entity by ID.
func (r *userProfileFieldRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.UserProfileField{
		Model: models.Model{
			ID: id,
		},
	}).Error
}
