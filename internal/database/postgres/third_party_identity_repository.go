package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// thirdPartyIdentityRepository stores and controls Users in the database.
type thirdPartyIdentityRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewThirdPartyIdentityRepository creates and returns a new PasswordResetRepository.
func NewThirdPartyIdentityRepository(db *gorm.DB, logger *zerolog.Logger) models.ThirdPartyIdentityRepository {
	return &thirdPartyIdentityRepository{db, logger}
}

// FindByID finds a single ThirdPartyIdentity by ID.
func (r *thirdPartyIdentityRepository) FindByID(id int64) (*models.ThirdPartyIdentity, error) {
	var tpi models.ThirdPartyIdentity
	if err := r.db.First(&tpi, id).Error; err != nil {
		return &tpi, err
	}
	return &tpi, nil
}

// FindByUserID finds a single ThirdPartyIdentity by UserID.
func (r *thirdPartyIdentityRepository) FindByUserID(userID int64) (*models.ThirdPartyIdentity, error) {
	var tpi models.ThirdPartyIdentity
	if err := r.db.Where("user_id = ?", userID).First(&tpi).Error; err != nil {
		return &tpi, err
	}

	return &tpi, nil
}

// FindUserIdentityByServiceName finds a single entity by user ID and service name.
func (r *thirdPartyIdentityRepository) FindUserIdentityByServiceName(userID int64, serviceName string) (*models.ThirdPartyIdentity, error) {
	var tpi models.ThirdPartyIdentity
	if err := r.db.Where("user_id = ? AND third_party_service_name = ?", userID, serviceName).First(&tpi).Error; err != nil {
		return &tpi, err
	}

	return &tpi, nil
}

// Create creates a new ThirdPartyIdentity.
func (r *thirdPartyIdentityRepository) Create(tpi models.ThirdPartyIdentity) error {
	return r.db.Create(&tpi).Error
}

// Update updates a ThirdPartyIdentity with the ID in the provided ThirdPartyIdentity.
func (r *thirdPartyIdentityRepository) Update(tpi models.ThirdPartyIdentity) error {
	return r.db.Model(&models.ThirdPartyIdentity{}).Updates(tpi).Error
}

// DeleteByID deletes a ThirdPartyIdentity by ID.
func (r *thirdPartyIdentityRepository) DeleteByID(id int64) error {
	return r.db.Delete(&models.ThirdPartyIdentity{
		Model: models.Model{
			ID: id,
		},
	}).Error
}
