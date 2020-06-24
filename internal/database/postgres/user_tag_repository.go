package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// userTagRepository represents the postgres implementation of the UserTagRepository.
type userTagRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewUserTagRepository creates and returns a new UserTagRepository with the provided parameters.
func NewUserTagRepository(db *gorm.DB, logger *zerolog.Logger) models.UserTagRepository {
	return &userTagRepository{
		db,
		logger,
	}
}

// FindByID finds a single entity by ID.
func (r *userTagRepository) FindByID(id int64) (*models.UserTag, error) {
	var userTag models.UserTag
	if err := r.db.First(&userTag, id).Error; err != nil {
		return &userTag, err
	}
	return &userTag, nil
}

// FindByUserID finds entities by UserID.
func (r *userTagRepository) FindByUserID(userID int64) ([]models.UserTag, error) {
	var userTags []models.UserTag
	if err := r.db.Where("user_id = ?", userID).Find(&userTags).Error; err != nil {
		return userTags, err
	}
	return userTags, nil
}

// Create creates a new entity.
func (r *userTagRepository) Create(userTag models.UserTag) error {
	return r.db.Create(&userTag).Error
}

// Update updates an entity with the ID in the provided entity.
func (r *userTagRepository) Update(userTag models.UserTag) error {
	return r.db.Model(&models.UserTag{}).Updates(userTag).Error
}

// DeleteByID deletes an entity by ID.
func (r *userTagRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.UserTag{
		Model: models.Model{
			ID: id,
		},
	}).Error
}
