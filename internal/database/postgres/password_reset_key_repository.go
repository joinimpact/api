package postgres

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// repository stores and controls Users in the database.
type passwordResetKeyRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewPasswordResetRepository creates and returns a new PasswordResetRepository.
func NewPasswordResetRepository(db *gorm.DB, logger *zerolog.Logger) models.PasswordResetKeyRepository {
	return &passwordResetKeyRepository{db, logger}
}

// FindByID finds a single PasswordResetKey by ID.
func (r *passwordResetKeyRepository) FindByID(id int64) (*models.PasswordResetKey, error) {
	var pwk models.PasswordResetKey
	if err := r.db.First(&pwk, id).Error; err != nil {
		return &pwk, err
	}
	return &pwk, nil
}

// FindByKey finds a single PasswordResetKey by Key.
func (r *passwordResetKeyRepository) FindByKey(key string) (*models.PasswordResetKey, error) {
	var pwk models.PasswordResetKey
	if err := r.db.Where("key = ?", key).First(&pwk).Error; err != nil {
		return &pwk, err
	}

	// Compare expired time.
	if pwk.ExpiresAt.Sub(time.Now().UTC()) <= 0 {
		return nil, errors.New("expired")
	}

	return &pwk, nil
}

// Create creates a new PasswordResetKey.
func (r *passwordResetKeyRepository) Create(pwk models.PasswordResetKey) error {
	return r.db.Create(&pwk).Error
}

// Update updates a PasswordResetKey with the ID in the provided PasswordResetKey.
func (r *passwordResetKeyRepository) Update(pwk models.PasswordResetKey) error {
	return r.db.Model(&models.PasswordResetKey{}).Updates(pwk).Error
}

// DeleteByID deletes a PasswordResetKey by ID.
func (r *passwordResetKeyRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.PasswordResetKey{
		Model: models.Model{
			ID: id,
		},
	}).Error
}
